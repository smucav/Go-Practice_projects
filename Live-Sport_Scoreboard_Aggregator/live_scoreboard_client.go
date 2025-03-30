// client module to handle businesses from client

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// dial to connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Server Down")
		return
	}

	defer conn.Close()

	go func() {
		// used to receive message from the server
		reader := bufio.NewReader(conn)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Server Disconnected")
				os.Exit(0)
			}
			fmt.Printf("%s\n> ", input)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		// use to send message to the server
		fmt.Printf("> ")
		input, _ := inputReader.ReadString('\n')
		msg := strings.TrimSpace(input)

		// send message to server
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Send Error", err)
			break
		}

		if msg == "EXIT" {
			break
		}

	}
}
