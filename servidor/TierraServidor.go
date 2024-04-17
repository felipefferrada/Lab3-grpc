package main

import (
	"context"
	"log"
	"net"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

// Define una estructura que implementa el servicio gRPC
type server struct {
	pb.UnimplementedRegionalServerServer
}

// Implementa el método ReceiveMessage del servicio gRPC
func (s *server) ReceiveMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	log.Printf("Received message: %v", in.Content)
	return &pb.Response{Message: "Message received"}, nil
}

func main() {
	// Define el puerto en el que el servidor escuchará las conexiones
	port := ":50051"

	// Crea un nuevo servidor gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}
	s := grpc.NewServer()

	// Registra el servicio en el servidor gRPC
	pb.RegisterRegionalServerServer(s, &server{})

	log.Printf("Server listening on port %v", port)

	// Inicia el servidor gRPC
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
