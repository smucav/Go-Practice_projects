package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to listen for server messages
	go listenForServerMessages(conn)

	// Send bids to the server from the terminal
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter bid (or type 'exit' to quit): ")
		input, _ := inputReader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error sending bid:", err)
			break
		}
	}
}

func listenForServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by server")
			os.Exit(0)
		}

		message = strings.TrimSpace(message)

		if message == "EXIT" {
			fmt.Println("\nServer is shutting down... Exiting client.")
			os.Exit(0)
		}

		fmt.Println("Server:", message)
	}
}
