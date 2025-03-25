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
		fmt.Println("Can't connect to server")
		return
	}

	reader := bufio.NewReader(conn)

	go func() {
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading")
				return
			}
			msg := strings.TrimSpace(input)
			fmt.Println(msg)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		msg, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error ", err)
			continue
		}

		msg = strings.TrimSpace(msg)

		_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Send Error", err)
			continue
		}
	}
}
