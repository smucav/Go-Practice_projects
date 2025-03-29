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
		fmt.Println("Server Down")
		return
	}

	defer conn.Close()

	go func() {
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
		fmt.Printf("> ")
		input, _ := inputReader.ReadString('\n')
		msg := strings.TrimSpace(input)

		_, err := conn.Write([]byte(msg + "\n"))
		if msg == "EXIT" {
			break
		}

		if err != nil {
			fmt.Println("Send Error", err)
		}
	}
}
