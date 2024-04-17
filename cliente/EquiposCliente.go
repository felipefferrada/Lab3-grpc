package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectarse al servidor: %v", err)
	}
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)

	msg := &pb.Message{Text: "Hola desde el cliente"}
	resp, err := c.SendMessage(context.Background(), msg)
	if err != nil {
		log.Fatalf("Error al enviar mensaje: %v", err)
	}

	fmt.Println("Respuesta del servidor:", resp.GetText())
}
