package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

var cantAT int = 0
var cantMP int = 0

type server struct {
	pb.UnimplementedChatServiceServer
}

func solicitarM(ID int, AT int, MP int) bool {
	if AT <= cantAT && MP <= cantMP { //es posible entregar
		cantAT -= AT
		cantMP -= MP
		//printeo solicitud
		return true
	} else {
		//printeo
		return false
	}
}

func abastecerAlmacen() {
	for {
		time.Sleep(5 * time.Second)

		if cantAT < 50 {
			cantAT += 10
		} else {
			cantAT = 50
		}

		if cantMP < 20 {
			cantMP += 5
		} else {
			cantMP = 20
		}
	}
}

func (s *server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	msgEquipo := msg.GetText()
	solicitud := strings.Split(msgEquipo, ",")
	idEquipo, err0 := strconv.Atoi(solicitud[0])
	solAT, err1 := strconv.Atoi(solicitud[1])
	solMP, err2 := strconv.Atoi(solicitud[2])

	if err0 != nil || err1 != nil || err2 != nil {
		log.Fatalf("%v", err0)
	}

	solicitarM(idEquipo, solAT, solMP)

	log.Printf("Mensaje recibido: %s", msgEquipo)

	// Aquí puedes procesar el mensaje según sea necesario
	// En este ejemplo, si el mensaje contiene "Terminar", el servidor enviará una respuesta para que la hebra del cliente se detenga
	if msg.GetText() == "Terminar" {
		return &pb.Message{Text: "Terminar"}, nil
	}

	// Si no es un mensaje para terminar, el servidor simplemente devuelve un mensaje de confirmación
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

	//generacion de municion en almacen
	go abastecerAlmacen()
}
