// client to listen and send to the broadcast server
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// connect the client to the broadcast server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}
	// close the connection after finishing
	defer conn.Close()

	// listen from the server message
	go func() {
		// read from the connection
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

	// start reader to read from the standar input
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, _ := inputReader.ReadString('\n')

		// clean unnecessary characters
		msg := strings.TrimSpace(input)

		// send to the server
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
