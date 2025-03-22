// project to practice goroutines, channels, select, timeouts,
// structs, methods, mutexes, networking, tickers, and rate limiting
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// client have conn to connect with broadcast server which helps to send and listen from/to broadcast server
type Client struct {
	conn          net.Conn
	mu            sync.Mutex
	burstyLimiter chan time.Time // rate limiter
	id            int
}

// message structure
type Message struct {
	SourceId int
	Content  string
}

// broadcast server have broadcastCh to listen from clients to broadcast and list of clients
// and assign id for clients
type BroadcastServer struct {
	broadcastCh chan Message
	clients     map[int]*Client
	mu          sync.Mutex
	nextID      int
}

// initialize broadcast server
func NewBroadcastServer() *BroadcastServer {
	return &BroadcastServer{
		clients:     make(map[int]*Client),
		nextID:      1,
		broadcastCh: make(chan Message, 10),
	}
}

func (bs *BroadcastServer) Run() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Server Failed")
		return
	}

	defer ln.Close()

	fmt.Println("Server started on :8080")

	// listen for client
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			bs.handleNewClient(conn)
		}
	}()

	// listen for message
	for msg := range bs.broadcastCh {
		bs.mu.Lock()
		// broadcast to each clients the received message
		for _, client := range bs.clients {
			client.mu.Lock()
			_, err := fmt.Fprintf(client.conn, "\nClient %d: %s\n", msg.SourceId, msg.Content)
			if err != nil {
				fmt.Println("Broadcase Error for client %d\n", client.id)
				client.mu.Unlock()
				continue
			}
			client.mu.Unlock()
		}
		bs.mu.Unlock()
	}

}

// handle new client
func (bs *BroadcastServer) handleNewClient(conn net.Conn) {
	bs.mu.Lock()
	clientID := bs.nextID
	bs.nextID++

	burstyLimiter := make(chan time.Time, 3)
	for i := 1; i <= 3; i++ {
		burstyLimiter <- time.Now()
	}

	// after the first 3 bursty then next message comes after 500 milliseconds
	go func() {
		for t := range time.Tick(500 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	client := &Client{
		conn:          conn,
		burstyLimiter: burstyLimiter,
		id:            clientID,
	}
	// register the client in broadcast server
	bs.clients[clientID] = client
	bs.mu.Unlock()

	fmt.Printf("Client %d connected\n", client.id)
	go bs.handleClient(client)
}

func (bs *BroadcastServer) handleClient(client *Client) {
	// after clients leave triggered function
	defer func() {
		bs.mu.Lock()
		delete(bs.clients, client.id) // remove the client so it will not accept any broadcasted message after leaving
		client.conn.Close()
		bs.mu.Unlock()
		fmt.Printf("client %d disconnected\n", client.id)
	}()

	// read from the client through it's conn
	reader := bufio.NewReader(client.conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return
			p
		}

		msg := strings.TrimSpace(input)

		if msg == "exit" {
			return
		}
		// after the first 3 message it will allow client to send message after 500 Millisecond
		// to avoid any busy traffic
		<-client.burstyLimiter
		fmt.Printf("Client %d sent message: %s\n", client.id, msg)

		bs.broadcastCh <- Message{SourceId: client.id, Content: msg}
	}
}

func main() {
	server := NewBroadcastServer()
	//start the server
	server.Run()
}
