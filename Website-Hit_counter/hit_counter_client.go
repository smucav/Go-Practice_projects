package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Connection err", err)
		return
	}

	// accept any message from the server using goroutine because
	// not to block at this level so it will be interactive server can send and also clients too
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Server Disconnected", err)
				break
			}

			fmt.Println(message)
		}
	}()

	// read from input to send to the server
	inputReader := bufio.NewReader(os.Stdin)

	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading ")
			break
		}

		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Send Error", err)
			continue
		}
	}
}
