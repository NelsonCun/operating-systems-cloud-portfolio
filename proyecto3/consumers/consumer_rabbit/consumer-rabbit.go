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
	// --- Conexión a Valkey/Redis ---
	addr := os.Getenv("VALKEY_SERVICE_URL")
	if addr == "" {
		addr = "valkey:6379" // Service de Redis/Valkey en Kubernetes
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Error conectando a Valkey:", err)
	} else {
		fmt.Println("✅ Conexión exitosa a Valkey:", pong)
	}

	// --- Conexión a RabbitMQ ---
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		log.Fatal("RABBITMQ_URL is required")
	}

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatal("Error conectando a RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("clima", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📥 Esperando mensajes de RabbitMQ...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Println("❌ Error parseando JSON:", err)
				continue
			}

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

			saveReading(rdb, municipio, temp, hum, weather)
			time.Sleep(1 * time.Second)
		}
	}()

	<-forever
}

// --- Guarda el dato y actualiza estadísticas ---
func saveReading(rdb *redis.Client, municipality string, temp int, hum int, weather string) {
	ts := time.Now().Unix()
	timeStr := time.Unix(ts, 0).Format("2006-01-02 15:04:05")

	// Último valor
	rdb.HSet(ctx, fmt.Sprintf("municipality:%s", municipality),
		"name", municipality,
		"temperature", temp,
		"humidity", hum,
		"weather", weather,
		"last_update", timeStr,
	)
	rdb.Expire(ctx, fmt.Sprintf("municipality:%s", municipality), 1*time.Hour)

	// Histórico
	rdb.LPush(ctx, fmt.Sprintf("temperature:%s", municipality), temp)
	rdb.LPush(ctx, fmt.Sprintf("humidity:%s", municipality), hum)
	rdb.LTrim(ctx, fmt.Sprintf("temperature:%s", municipality), 0, 999)
	rdb.LTrim(ctx, fmt.Sprintf("humidity:%s", municipality), 0, 999)

	// Series de tiempo
	rdb.ZAdd(ctx, fmt.Sprintf("temperature_prom:%s", municipality), redis.Z{
		Score:  float64(temp),
		Member: fmt.Sprintf("%d", ts),
	})
	rdb.ZAdd(ctx, fmt.Sprintf("humidity_prom:%s", municipality), redis.Z{
		Score:  float64(hum),
		Member: fmt.Sprintf("%d", ts),
	})

	// --- Promedios acumulados ---
	accTempKey := fmt.Sprintf("municipality:%s:temp_sum", municipality)
	accHumKey := fmt.Sprintf("municipality:%s:hum_sum", municipality)
	countKey := fmt.Sprintf("municipality:%s:count", municipality)

	rdb.IncrByFloat(ctx, accTempKey, float64(temp))
	rdb.IncrByFloat(ctx, accHumKey, float64(hum))
	rdb.Incr(ctx, countKey)

	sumTemp, _ := rdb.Get(ctx, accTempKey).Float64()
	sumHum, _ := rdb.Get(ctx, accHumKey).Float64()
	count, _ := rdb.Get(ctx, countKey).Float64()

	if count > 0 {
		promTemp := sumTemp / count
		promHum := sumHum / count

		rdb.ZAdd(ctx, "total:temperature_prom", redis.Z{Score: promTemp, Member: municipality})
		rdb.ZAdd(ctx, "total:humidity_prom", redis.Z{Score: promHum, Member: municipality})
	}

	// --- Mínimos y máximos totales ---
	updateMinMax(rdb, "total:temperature:min", "total:temperature:max", float64(temp))
	updateMinMax(rdb, "total:humidity:min", "total:humidity:max", float64(hum))

	// --- Conteo de condiciones ---
	rdb.ZIncrBy(ctx, fmt.Sprintf("conteo_clima:%s", municipality), 1, weather)
	rdb.ZIncrBy(ctx, "conteo_clima:total", 1, weather)

	// Clima más y menos común
	top, _ := rdb.ZRevRange(ctx, fmt.Sprintf("conteo_clima:%s", municipality), 0, 0).Result()
	if len(top) > 0 {
		rdb.HSet(ctx, fmt.Sprintf("municipality:%s", municipality), "clima_mas_comun", top[0])
	}
	least, _ := rdb.ZRange(ctx, fmt.Sprintf("conteo_clima:%s", municipality), 0, 0).Result()
	if len(least) > 0 {
		rdb.HSet(ctx, fmt.Sprintf("municipality:%s", municipality), "clima_menos_comun", least[0])
	}

	fmt.Printf("💾 Guardado en Valkey: %s -> Temp:%d, Hum:%d, Weather:%s\n", municipality, temp, hum, weather)
}

// --- Actualiza valores totales min/max ---
func updateMinMax(rdb *redis.Client, keyMin, keyMax string, value float64) {
	currentMin, err := rdb.Get(ctx, keyMin).Float64()
	if err == redis.Nil || value < currentMin {
		rdb.Set(ctx, keyMin, value, 0)
	}
	currentMax, err := rdb.Get(ctx, keyMax).Float64()
	if err == redis.Nil || value > currentMax {
		rdb.Set(ctx, keyMax, value, 0)
	}
}
