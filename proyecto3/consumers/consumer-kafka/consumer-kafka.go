package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var ctx = context.Background()

func main() {
	// ----------------------
	// Configurar cliente Valkey
	// ----------------------
	addr := func() string {
		if v := os.Getenv("VALKEY_SERVICE_URL"); v != "" {
			return v
		}
		return "localhost:6379"
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Error conectando a Valkey:", err)
	} else {
		fmt.Println("Conexión exitosa a Valkey:", pong)
	}

	// ----------------------
	// Conexión Kafka
	// ----------------------
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "clima",
		GroupID: "clima-consumer-group",
	})

	fmt.Println("📥 Esperando mensajes de Kafka...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("❌ Error leyendo mensaje:", err)
			continue
		}

		fmt.Printf("✅ [Kafka] Mensaje recibido: %s\n", string(m.Value))

		// Parsear JSON
		var data map[string]interface{}
		if err := json.Unmarshal(m.Value, &data); err != nil {
			log.Println("❌ Error parseando JSON:", err)
			continue
		}

		// Validar y extraer campos
		municipality, ok := data["municipality"].(string)
		if !ok {
			log.Println("❌ Campo 'municipality' no encontrado o no es string")
			continue
		}

		temperature, ok := data["temperature"].(float64)
		if !ok {
			log.Println("❌ Campo 'temperature' no encontrado o no es número")
			continue
		}

		humidity, ok := data["humidity"].(float64)
		if !ok {
			log.Println("❌ Campo 'humidity' no encontrado o no es número")
			continue
		}

		weather, ok := data["weather"].(string)
		if !ok {
			log.Println("❌ Campo 'weather' no encontrado o no es string")
			continue
		}

		// Guardar en Valkey
		saveReading(rdb, municipality, int(temperature), int(humidity), weather)

		time.Sleep(1 * time.Second)
	}

}

// ----------------------
// Función para guardar lectura en Valkey
// ----------------------
func saveReading(rdb *redis.Client, municipality string, temp int, hum int, weather string) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	key := fmt.Sprintf("municipality:%s", municipality)

	if err := rdb.HSet(ctx,
		key,
		"name", municipality,
		"temperature", temp,
		"humidity", hum,
		"weather", weather,
		"last_update", ts,
	).Err(); err != nil {
		log.Println("❌ Error guardando en Valkey:", err)
		return
	}

	rdb.Expire(ctx, key, 1*time.Hour)
	fmt.Printf("💾 Guardado en Valkey: %s -> Temp:%d, Hum:%d, Weather:%s\n", municipality, temp, hum, weather)
}
