package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn          net.Conn
	mu            sync.Mutex
	burstyLimiter chan time.Time
	id            int
}

type Message struct {
	SourceId int
	Content  string
}

type BroadcastServer struct {
	broadcastCh chan Message
	clients     map[int]*Client
	mu          sync.Mutex
	nextID      int
}

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

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			bs.handleNewClient(conn)
		}
	}()

	for {
		for msg := range bs.broadcastCh {
			bs.mu.Lock()
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
}

func (bs *BroadcastServer) handleNewClient(conn net.Conn) {
	bs.mu.Lock()
	clientID := bs.nextID
	bs.nextID++

	burstyLimiter := make(chan time.Time, 3)
	for i := 1; i <= 3; i++ {
		burstyLimiter <- time.Now()
	}

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
	bs.clients[clientID] = client
	bs.mu.Unlock()

	fmt.Printf("Client %d connected\n", client.id)
	go bs.handleClient(client)
}

func (bs *BroadcastServer) handleClient(client *Client) {
	defer func() {
		bs.mu.Lock()
		delete(bs.clients, client.id)
		client.conn.Close()
		bs.mu.Unlock()
		fmt.Printf("client %d disconnected\n", client.id)
	}()

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
		<-client.burstyLimiter
		fmt.Printf("Client %d sent message: %s\n", client.id, msg)

		bs.broadcastCh <- Message{SourceId: client.id, Content: msg}
	}
}

func main() {
	server := NewBroadcastServer()
	server.Run()
}
