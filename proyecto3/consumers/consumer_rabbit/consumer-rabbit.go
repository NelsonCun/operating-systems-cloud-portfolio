package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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
	// Conexión RabbitMQ
	// ----------------------
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("❌ Error conectando a RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("❌ Error abriendo canal:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("clima", false, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error declarando cola:", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes:", err)
	}

	fmt.Println("📥 Esperando mensajes de RabbitMQ...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("✅ [RabbitMQ] Mensaje recibido: %s\n", d.Body)

			// Parsear mensaje JSON en un mapa
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Println("❌ Error parseando JSON:", err)
				continue
			}

			// Validar y convertir cada campo de forma segura
			municipio, ok := data["municipality"].(string)
			if !ok {
				log.Println("❌ Clave 'municipality' faltante o no es string")
				continue
			}

			tempFloat, ok := data["temperature"].(float64)
			if !ok {
				log.Println("❌ Clave 'temperature' faltante o no es float64")
				continue
			}
			temp := int(tempFloat)

			humFloat, ok := data["humidity"].(float64)
			if !ok {
				log.Println("❌ Clave 'humidity' faltante o no es float64")
				continue
			}
			hum := int(humFloat)

			weather, ok := data["weather"].(string)
			if !ok {
				log.Println("❌ Clave 'weather' faltante o no es string")
				continue
			}

			// Guardar en Valkey
			saveReading(rdb, municipio, temp, hum, weather)

			time.Sleep(1 * time.Second)
		}
	}()

	<-forever
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
