package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/felipefferrada/Lab3-grpc/proto"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:50051"
)

func sendMessage(client pb.ChatServiceClient, wg *sync.WaitGroup, id int, message string, responseCh chan<- bool) {
	defer wg.Done()
	fmt.Printf(message + "\n")

	for {
		msg := &pb.Message{Text: message}
		resp, err := client.SendMessage(context.Background(), msg)

		if err != nil {
			log.Printf("Error al enviar mensaje: %v", err)
			time.Sleep(time.Second) // Reintentar después de un descanso.
			continue
		}

		//message contiene la solicitud (ID equipo, at, mp)
		listaMensaje := strings.Split(message, ", ")
		idEquipo, err0 := strconv.Atoi(listaMensaje[0])
		solAT, err1 := strconv.Atoi(listaMensaje[1])
		solMP, err2 := strconv.Atoi(listaMensaje[2])

		if err0 != nil || err1 != nil || err2 != nil {
			log.Fatalf("%v", err0)
		}
		//el print con la respuesta y la solicitud
		//fmt.Printf("Hebra %d - Respuesta del servidor: %s\n", id, resp.GetText())

		// Dependiendo de la respuesta, decide si enviar otro mensaje o terminar la hebra
		if resp.GetText() == "true" || resp.GetText() == "1" {
			responseCh <- true // Indica que la hebra debe terminar
			fmt.Printf("Equipo %d Solicitando %d AT y %d MP ; Resolucion: -- APROBADA -- ;\nConquista Exitosa!, cerrando comunicacion", idEquipo, solAT, solMP)
			return
		} else {
			fmt.Printf("Equipo %d Solicitando %d AT y %d MP ; Resolucion: -- DENEGADA -- ;\nReintentando en 3 segs...", idEquipo, solAT, solMP)

		}

		time.Sleep(3 * time.Second) // Espera 3 segundos antes de enviar otro mensaje
	}
}

func main() {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectarse al servidor: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	var wg sync.WaitGroup
	responseCh := make(chan bool)

	// Inicia cuatro hebras concurrentes
	//el mensaje enviado es tipo (id Equipo, cantAT, cantMP)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go sendMessage(client, &wg, i+1, fmt.Sprintf("%d, %d, %d", i+1, rand.Intn(30-20+1)+20, rand.Intn(15-10+1)+10), responseCh)
	}

	// Espera a que todas las hebras terminen
	go func() {
		wg.Wait()
		close(responseCh)
	}()

	// Espera respuestas de las hebras para determinar si deben enviar más mensajes
	for resp := range responseCh {
		if resp {
			fmt.Println("Hebra terminada.")
		} else {
			fmt.Println("Hebra enviando otro mensaje.")
		}
	}
}
