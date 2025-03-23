package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Connection err", err)
		return
	}

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
