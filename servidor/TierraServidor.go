package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

var (
	cantAT int = 0
	cantMP int = 0
	mu     sync.Mutex
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func solicitarM(ID int, AT int, MP int) bool {
	if AT <= cantAT && MP <= cantMP { //es posible entregar
		cantAT -= AT
		cantMP -= MP
		fmt.Printf("Recepcion de solicitud desde equipo %d, %d AT y %d MP -- APROBADA --\n", ID, AT, MP)
		fmt.Printf("AT EN SISTEMA: %d ; MP EN SISTEMA: %d \n", cantAT, cantMP)
		//printeo solicitud
		return true
	} else {
		fmt.Printf("Recepcion de solicitud desde equipo %d, %d AT y %d MP -- DENEGADA --", ID, AT, MP)
		fmt.Printf("AT EN SISTEMA: %d ; MP EN SISTEMA: %d \n", cantAT, cantMP)
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
	mu.Lock()
	defer mu.Unlock()
	msgEquipo := msg.GetText()
	solicitud := strings.Split(msgEquipo, ", ")
	idEquipo, err0 := strconv.Atoi(solicitud[0])
	solAT, err1 := strconv.Atoi(solicitud[1])
	solMP, err2 := strconv.Atoi(solicitud[2])

	if err0 != nil || err1 != nil || err2 != nil {
		log.Fatalf("%v", err0)
	}

	respuesta := solicitarM(idEquipo, solAT, solMP)
	//log.Printf("Mensaje recibido: %s", msgEquipo)

	// Aquí puedes procesar el mensaje según sea necesario
	if respuesta {
		return &pb.Message{Text: "true"}, nil
	} else {
		// Si no es un mensaje para terminar, el servidor simplemente devuelve un mensaje de confirmación
		return &pb.Message{Text: "false"}, nil
	}

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})

	//generacion de municion en almacen
	go abastecerAlmacen()

	fmt.Println("Servidor gRPC iniciado en el puerto :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}

}
