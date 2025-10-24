package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

type Clima struct {
	Municipality string `json:"municipality"`
	Temperature  int    `json:"temperature"`
	Humidity     int    `json:"humidity"`
	Weather      string `json:"weather"`
}

var municipalities = []string{"Mixco", "Guatemala", "Amatitlan", "Chinautla"}
var weathers = []string{"sunny", "cloudy", "rainy", "foggy"}

func generarClima() Clima {
	return Clima{
		Municipality: municipalities[rand.Intn(len(municipalities))],
		Temperature:  15 + rand.Intn(15),
		Humidity:     50 + rand.Intn(50),
		Weather:      weathers[rand.Intn(len(weathers))],
	}
}

// Enviar a Kafka
func enviarKafka(msg Clima, kafkaBroker string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, "clima", 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, _ := json.Marshal(msg)
	_, err = conn.WriteMessages(kafka.Message{Value: data})
	return err
}

// Enviar a RabbitMQ
func enviarRabbit(msg Clima, rabbitmqURL string) error {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("clima", false, false, false, false, nil)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(msg)
	return ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Configuración desde variables de entorno
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	fmt.Printf("🚀 Enviando datos de clima...\n")
	fmt.Printf("📡 Kafka broker: %s\n", kafkaBroker)
	fmt.Printf("🐰 RabbitMQ URL: %s\n", rabbitmqURL)

	for {
		clima := generarClima()
		data, _ := json.MarshalIndent(clima, "", "  ")
		fmt.Println("📦 Nuevo dato:", string(data))

		if rand.Intn(2) == 0 {
			fmt.Println("➡️  Enviando a Kafka")
			if err := enviarKafka(clima, kafkaBroker); err != nil {
				log.Println("❌ Error Kafka:", err)
			} else {
				fmt.Println("✅ Enviado a Kafka exitosamente")
			}
		} else {
			fmt.Println("➡️  Enviando a RabbitMQ")
			if err := enviarRabbit(clima, rabbitmqURL); err != nil {
				log.Println("❌ Error RabbitMQ:", err)
			} else {
				fmt.Println("✅ Enviado a RabbitMQ exitosamente")
			}
		}

		time.Sleep(3 * time.Second)
	}
}

// Comando para consumir desde Kafka:
// kubectl exec -it <kafka-pod-name> -n clima-app -- kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic clima --from-beginning
