package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	pb "server_go/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherTweetServiceServer
}

var (
	// Se inicializa con valores por defecto
	kafkaWriterAddr  = "localhost:6001"
	rabbitWriterAddr = "localhost:6002"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Sobrescribir si existen variables de entorno
	if env := os.Getenv("KAFKA_BROKER"); env != "" {
		kafkaWriterAddr = env
	}
	if env := os.Getenv("RABBITMQ_URL"); env != "" {
		rabbitWriterAddr = env
	}
}

// SendTweet decide aleatoriamente a cuál writer enviar
func (s *server) SendTweet(ctx context.Context, req *pb.WeatherTweetRequest) (*pb.WeatherTweetResponse, error) {
	fmt.Printf("📩 gRPC recibió: %+v\n", req)

	var writerAddr string
	if rand.Intn(2) == 0 {
		writerAddr = kafkaWriterAddr
		fmt.Println("➡️  Enviando a Kafka Writer")
	} else {
		writerAddr = rabbitWriterAddr
		fmt.Println("➡️  Enviando a RabbitMQ Writer")
	}

	// Conectarse al writer por gRPC
	conn, err := grpc.Dial(writerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("❌ Error conectando al writer:", err)
		return &pb.WeatherTweetResponse{Status: "Error conectando al writer"}, err
	}
	defer conn.Close()

	client := pb.NewWeatherTweetServiceClient(conn)
	resp, err := client.SendTweet(ctx, req)
	if err != nil {
		log.Println("❌ Error enviando al writer:", err)
		return &pb.WeatherTweetResponse{Status: "Error enviando al writer"}, err
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Error escuchando: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherTweetServiceServer(grpcServer, &server{})

	fmt.Println("🚀 gRPC Server escuchando en :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Error en Serve: %v", err)
	}
}
