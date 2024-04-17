package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Printf("Mensaje recibido: %s", msg.GetText())
	return &pb.Message{Text: "Mensaje recibido"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})

	fmt.Println("Servidor gRPC iniciado en el puerto :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
