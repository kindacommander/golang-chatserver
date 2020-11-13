package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client struct {
	ch   chan<- string // Outbound channel
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // All client's incoming messages
)

func broadcaster() {
	clients := make(map[client]bool) // All connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// outbound channels for clients.
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			clients[cli] = true

			var currentClients []string
			for c := range clients {
				currentClients = append(currentClients, c.name)
			}
			cli.ch <- fmt.Sprintf("%d users online:", len(currentClients))

			for _, c := range currentClients {
				cli.ch <- fmt.Sprintf("%s", c)
			}
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // Client's outbound channel
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " connected"
	entering <- client{ch, who}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// Ignoring potential input.Err()

	leaving <- client{ch, who}
	messages <- who + " disconnected"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // Ignoring network errors
	}
}
