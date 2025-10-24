package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "writer_kafka/proto"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// ---------------------------
// Struct para enviar a Kafka
// ---------------------------
type KafkaClima struct {
	Municipality string `json:"municipality"`
	Temperature  int32  `json:"temperature"`
	Humidity     int32  `json:"humidity"`
	Weather      string `json:"weather"`
}

// ---------------------------
// Servidor gRPC
// ---------------------------
type server struct {
	pb.UnimplementedWeatherTweetServiceServer
}

var kafkaBroker string

func init() {
	kafkaBroker = os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
}

// ---------------------------
// Función para enviar a Kafka
// ---------------------------
func enviarKafka(msg KafkaClima) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, "clima", 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = conn.WriteMessages(kafka.Message{Value: data})
	return err
}

// ---------------------------
// Método gRPC
// ---------------------------
func (s *server) SendTweet(ctx context.Context, req *pb.WeatherTweetRequest) (*pb.WeatherTweetResponse, error) {
	municipioNombre := req.Municipality.String()
	if municipioNombre == "" {
		municipioNombre = "Desconocido"
	}

	weatherNombre := req.Weather.String()

	kafkaMsg := KafkaClima{
		Municipality: municipioNombre,
		Temperature:  req.Temperature,
		Humidity:     req.Humidity,
		Weather:      weatherNombre,
	}

	fmt.Printf("📦 Kafka Writer recibió: %+v\n", kafkaMsg)
	if err := enviarKafka(kafkaMsg); err != nil {
		log.Println("❌ Error Kafka:", err)
		return &pb.WeatherTweetResponse{Status: "Error Kafka"}, err
	}

	fmt.Println("✅ Enviado a Kafka exitosamente")
	return &pb.WeatherTweetResponse{Status: "Enviado a Kafka"}, nil
}

// ---------------------------
// Función main
// ---------------------------
func main() {
	lis, err := net.Listen("tcp", ":6001")
	if err != nil {
		log.Fatalf("❌ Error escuchando Kafka Writer: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherTweetServiceServer(grpcServer, &server{})

	fmt.Println("🚀 Kafka Writer escuchando en :6001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Error en Serve Kafka Writer: %v", err)
	}
}
