package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "writer_rabbit/proto"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherTweetServiceServer
}

var rabbitmqURL string

func init() {
	rabbitmqURL = os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}
}

func enviarRabbit(msg *pb.WeatherTweetRequest) error {
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

	dataMap := map[string]interface{}{
		"municipality": msg.Municipality.String(),
		"temperature":  msg.Temperature,
		"humidity":     msg.Humidity,
		"weather":      msg.Weather.String(),
	}

	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (s *server) SendTweet(ctx context.Context, req *pb.WeatherTweetRequest) (*pb.WeatherTweetResponse, error) {
	fmt.Printf("📦 RabbitMQ Writer recibió: %+v\n", req)
	if err := enviarRabbit(req); err != nil {
		log.Println("❌ Error RabbitMQ:", err)
		return &pb.WeatherTweetResponse{Status: "Error RabbitMQ"}, err
	}
	fmt.Println("✅ Enviado a RabbitMQ exitosamente")
	return &pb.WeatherTweetResponse{Status: "Enviado a RabbitMQ"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":6002")
	if err != nil {
		log.Fatalf("❌ Error escuchando RabbitMQ Writer: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherTweetServiceServer(grpcServer, &server{})

	fmt.Println("🚀 RabbitMQ Writer escuchando en :6002")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Error en Serve RabbitMQ Writer: %v", err)
	}
}
