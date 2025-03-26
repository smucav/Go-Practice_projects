package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}

	defer conn.Close()

	closeCh := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(conn)
		for {
			select {
			case <-closeCh:
				return
			default:
				msg, err := reader.ReadString('.')
				if err != nil {
					fmt.Println("Error reading or Server Down:(")
					os.Exit(0)
				}
				fmt.Printf("%s \n> ", msg)
			}
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		msg, err := inputReader.ReadString('\n')

		msg = strings.TrimSpace(msg)

		_, err = conn.Write([]byte(msg + "\n"))

		if err != nil {
			fmt.Println("Send Error", err)
			continue
		}
		if msg == "exit" {
			close(closeCh)
			wg.Wait()
			return
		}

	}
}
