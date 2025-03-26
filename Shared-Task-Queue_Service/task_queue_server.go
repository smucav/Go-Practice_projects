package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type setTask struct {
	task string
}

type getTask struct{}

type client struct {
	conn    net.Conn
	setTask setTask
	getTask getTask
}

type TaskQueueServer struct {
	gets    chan client
	sets    chan client
	closeCh chan struct{}
	wg      sync.WaitGroup
}

func NewTaskQueueServer() *TaskQueueServer {
	return &TaskQueueServer{
		gets:    make(chan client),
		sets:    make(chan client),
		closeCh: make(chan struct{}),
	}
}

func (tqs *TaskQueueServer) Run() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Server Failed")
		return
	}
	defer ln.Close()
	fmt.Println("Server started on :8080")

	tqs.wg.Add(1)
	go func() {
		defer tqs.wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Connection error")
				return
			}
			tqs.wg.Add(1)
			go tqs.handleClient(conn)
		}
	}()

	go func() {
		var queue []string
		for {
			select {
			case getc := <-tqs.gets:
				if len(queue) == 0 {
					fmt.Fprintf(getc.conn, "There is No task have fun.")
				} else {
					fmt.Printf("Task %s Taken\n", queue[0])
					fmt.Fprintf(getc.conn, "Task -> "+queue[0]+".")
					queue = queue[1:]
				}

			case setc := <-tqs.sets:
				queue = append(queue, setc.setTask.task)
				fmt.Printf("Task %s Added\n", queue[len(queue)-1])
				fmt.Fprintf(setc.conn, "Task added successfully.")
			case <-tqs.closeCh:
				return
			}
		}
	}()

	tqs.wg.Wait()
	close(tqs.closeCh)
}

func (tqs *TaskQueueServer) handleClient(conn net.Conn) {
	defer tqs.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Welcome to task queue \nAdd task(Task: task you want to add) or get(NEXT).")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg := strings.TrimSpace(input)

		if msg == "exit" {
			fmt.Fprintf(conn, "Good bye :)\n")
			continue
		}

		if msg == "NEXT" {
			get := client{
				conn:    conn,
				getTask: getTask{},
			}
			tqs.gets <- get
			continue
		}

		if strings.HasPrefix(strings.ToLower(msg), "task") {
			task := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(msg), "task: "))
			if task == "" {
				fmt.Fprintf(conn, "Invalid Page\n")

				continue
			}
			set := client{
				conn:    conn,
				setTask: setTask{task: task},
			}
			tqs.sets <- set
			continue
		}
		fmt.Fprintf(conn, "Unknown command.")
	}
}

func main() {
	server := NewTaskQueueServer()
	server.Run()
}
