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
	conn, err := net.Dial("tcp", "localhost:8080") // connect to the server

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
		// read from the server
		reader := bufio.NewReader(conn)
		for {
			// used select here to listen to the client exit event and quit gracefully
			select {
			case <-closeCh: // listen to close channel signal
				return
			default:
				msg, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading or Server Down:(")
					os.Exit(0)
				}
				// print the message from the server
				fmt.Printf("%s> ", msg)
			}
		}
	}()

	// read from the client standard in
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		msg, err := inputReader.ReadString('\n')
		// remove trailing characters
		msg = strings.TrimSpace(msg)
		// send to the server
		_, err = conn.Write([]byte(msg + "\n"))

		if err != nil {
			fmt.Println("Send Error", err)
			continue
		}
		if msg == "exit" {
			// when client exit close the channel then help to signal the channel and exit the loop and go routine
			close(closeCh)
			wg.Wait()
			return
		}

	}
}
