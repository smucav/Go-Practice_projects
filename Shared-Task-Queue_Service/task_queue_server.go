package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type setTask struct {
	task   string
	respCh chan bool
}

type getTask struct {
	respCh chan string
}

type TaskQueueServer struct {
	gets    chan getTask
	sets    chan setTask
	closeCh chan struct{}
	wg      sync.WaitGroup
}

func NewTaskQueueServer() *TaskQueueServer {
	return &TaskQueueServer{
		gets:    make(chan getTask),
		sets:    make(chan setTask),
		closeCh: make(chan struct{}),
	}
}

func (tqs *TaskQueueServer) Run() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("connection error")
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error")
			return
		}
		tqs.wg.Add(1)
		go tqs.handleClient(conn)
	}

	tqs.wg.Add(1)
	go func() {
		defer tqs.wg.Done()
		queue := make([]string, 0)
		for {
			select {
			case get := <-tqs.gets:
				get.respCh <- queue[0]
				queue = queue[1:]
			case set := <-tqs.sets:
				queue = append(queue, set.task)
				set.respCh <- true
			case <-tqs.closeCh:
				return
			}
		}
	}()
}

func (tqs *TaskQueueServer) handleClient(conn net.Conn) {
	defer tqs.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg := strings.TrimSpace(input)

		if msg == "exit" {
			close(tqs.closeCh)
			return
		}

		if msg == "NEXT" {
			get := getTask{respCh: make(chan string)}
			tqs.gets <- get
		}

		if strings.HasPrefix(msg, "Task") {
			task := strings.TrimSpace(strings.TrimPrefix(msg, "Task: "))
			if task == "" {
				fmt.Fprintf(conn, "Invalid Page\n")

				continue
			}
			set := setTask{task: task, respCh: make(chan bool)}
			tqs.sets <- set
		}
	}
	tqs.wg.Add(1)
	go func() {
		defer tqs.wg.Done()
		for {
			select {
			case t := <-tqs.gets:
				task := <-t.respCh
				fmt.Fprintf(conn, task)
			case <-tqs.sets:
				fmt.Fprintf(conn, "Task set successfully\n")
			case <-tqs.closeCh:
				return
			}
		}
	}()
}

func main() {
	server := NewTaskQueueServer()
	server.Run()
}
