package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// ----------------------
	// Configurar cliente Valkey (Redis)
	// ----------------------
	addr := getEnv("VALKEY_SERVICE_URL", "localhost:6379")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Error conectando a Valkey:", err)
	} else {
		log.Println("Conexión exitosa a Valkey:", pong)
	}

	// ----------------------
	// Puerto configurable
	// ----------------------
	port := getEnv("PORT", "8081")

	// ----------------------
	// Endpoints HTTP
	// ----------------------
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		handleData(w, r, rdb)
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handleMetrics(w, r, rdb)
	})

	log.Printf("🚀 Servidor listo en puerto %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// ----------------------
// Handlers
// ----------------------
func handleData(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	municipality := r.URL.Query().Get("municipality")
	if municipality == "" {
		http.Error(w, "Se requiere parámetro 'municipality'", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("municipality:%s:history", municipality)
	vals, err := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		http.Error(w, "Error leyendo Valkey: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := []map[string]interface{}{}
	for _, z := range vals {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(z.Member.(string)), &data); err != nil {
			continue
		}
		data["time"] = int64(z.Score * 1000)
		result = append(result, data)
	}

	writeJSON(w, result)
}

func handleMetrics(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	metric := r.URL.Query().Get("metric")
	municipality := r.URL.Query().Get("municipality") // opcional

	if metric == "" {
		http.Error(w, "Se requiere parámetro 'metric'", http.StatusBadRequest)
		return
	}

	var keys []string
	if municipality != "" {
		keys = []string{fmt.Sprintf("municipality:%s:history", municipality)}
	} else {
		keys, _ = rdb.Keys(ctx, "municipality:*:history").Result()
	}

	if len(keys) == 0 {
		http.Error(w, "No se encontraron datos", http.StatusNotFound)
		return
	}

	var sumTemp, sumHum int
	var count int
	var maxTemp, maxHum int
	var minTemp, minHum int
	first := true
	weatherCount := make(map[string]int)

	for _, key := range keys {
		vals, _ := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
		for _, z := range vals {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(z.Member.(string)), &data); err != nil {
				continue
			}

			temp, _ := toInt(data["temperature"])
			hum, _ := toInt(data["humidity"])
			weather, _ := data["weather"].(string)

			sumTemp += temp
			sumHum += hum
			count++

			if first {
				maxTemp, maxHum = temp, hum
				minTemp, minHum = temp, hum
				first = false
			} else {
				if temp > maxTemp {
					maxTemp = temp
				}
				if temp < minTemp {
					minTemp = temp
				}
				if hum > maxHum {
					maxHum = hum
				}
				if hum < minHum {
					minHum = hum
				}
			}

			weatherCount[weather]++
		}
	}

	var value interface{}
	switch metric {
	case "humidity_avg":
		if count > 0 {
			value = float64(sumHum) / float64(count)
		} else {
			value = 0
		}
	case "temperature_avg":
		if count > 0 {
			value = float64(sumTemp) / float64(count)
		} else {
			value = 0
		}
	case "temperature_max":
		value = maxTemp
	case "humidity_max":
		value = maxHum
	case "temperature_min":
		value = minTemp
	case "humidity_min":
		value = minHum
	case "weather_count":
		value = weatherCount
	case "weather_most_common":
		value = mostCommon(weatherCount)
	case "weather_least_common":
		value = leastCommon(weatherCount)
	default:
		http.Error(w, "Métrica no soportada", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"metric":       metric,
		"municipality": municipality,
		"value":        value,
		"time":         time.Now().Unix() * 1000,
	}

	writeJSON(w, resp)
}

// ----------------------
// Helpers
// ----------------------
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func toInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case string:
		i, err := strconv.Atoi(v)
		return i, err == nil
	default:
		return 0, false
	}
}

func mostCommon(countMap map[string]int) string {
	maxCount := -1
	result := ""
	for k, v := range countMap {
		if v > maxCount {
			maxCount = v
			result = k
		}
	}
	return result
}

func leastCommon(countMap map[string]int) string {
	minCount := -1
	result := ""
	keys := make([]string, 0, len(countMap))
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := countMap[k]
		if minCount == -1 || v < minCount {
			minCount = v
			result = k
		}
	}
	return result
}
