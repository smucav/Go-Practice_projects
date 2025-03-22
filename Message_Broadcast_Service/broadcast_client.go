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
		fmt.Println("Connection error: ", err)
		return
	}
	defer conn.Close()

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Server Disconnected: ", err)
				os.Exit(0)
			}
			fmt.Println(msg)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, _ := inputReader.ReadString('\n')

		msg := strings.TrimSpace(input)

		_, err := conn.Write([]byte(msg + "\n"))

		if msg == "exit" {
			break
		}

		if err != nil {
			fmt.Println("Send err ", err)
			break
		}
	}
}
