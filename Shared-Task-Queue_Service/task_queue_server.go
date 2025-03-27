package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// struct each that will be submitted or take with id
type Task struct {
	ID      int
	Content string
}

// struct that contain task struct and also a response channel
// to verify whether the task submitted successfully or not
// think of it as carrier
type SubmitOp struct {
	task   Task
	respCh chan bool
}

// struct same as SubmitOp this will send the task to the client
// using a dedicated response channel
type NextOp struct {
	respCh chan *Task
}

// helps to check how many tasks left in the queue
type StatusOp struct {
	respCh chan int
}

// server which governs submit and take tasks from the queue
// and with waitgroup to help track of the clients
type TaskQueueServer struct {
	submitCh chan SubmitOp
	nextCh   chan NextOp
	statusCh chan StatusOp
	client   sync.WaitGroup
}

// Initilize the structs
func NewTaskQueueServer() *TaskQueueServer {
	return &TaskQueueServer{
		submitCh: make(chan SubmitOp),
		nextCh:   make(chan NextOp),
		statusCh: make(chan StatusOp),
	}
}

func (tqs *TaskQueueServer) Run() {
	// listen to any connection with this port
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Server error", err)
		return
	}
	// close listening after done
	defer ln.Close()

	fmt.Println("Task Queue Server started on :8080")

	// used goroutine here to not block accepting new clients
	// and to apply stateful goroutine concept
	go func() {
		// queue where the tasks will be stored
		var queue []Task
		// helps to track the id of tasks
		var taskId int
		for {
			// used select here to listen to incoming requests
			select {
			case submit := <-tqs.submitCh:
				// increase when new task submitted
				taskId++
				submit.task.ID = taskId
				queue = append(queue, submit.task)
				// send signal that the task is stored successfully
				submit.respCh <- true
			case next := <-tqs.nextCh:
				// return the first task as FIRST IN LAST OUT operation
				if len(queue) == 0 {
					next.respCh <- nil
				} else {
					task := queue[0]
					next.respCh <- &task
					queue = queue[1:]
				}
			case status := <-tqs.statusCh:
				// returns number of tasks left using channel
				status.respCh <- len(queue)
			}
		}
	}()

	for {
		// accept new clients
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Server Error", err)
			return
		}

		// initilize new go routine for each clients
		// to handle them concurrently
		tqs.client.Add(1)
		go tqs.handleClient(conn)
	}
}

func (tqs *TaskQueueServer) handleClient(conn net.Conn) {
	defer conn.Close()
	defer tqs.client.Done()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Welcome to Task Queue!\nAdd Task (Task: <task>) or get (NEXT), Status (status), 'exit' to quit:\n")

	for {
		// accept client message
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected\n", conn.RemoteAddr().String())
			return
		}

		// remove trailing characters
		msg := strings.TrimSpace(input)

		// check the message whether they comply with commands avaliable or not

		if msg == "exit" {
			fmt.Printf("Client %s exited\n", conn.RemoteAddr().String())
			fmt.Fprintf(conn, "Good bye :)\n")
			return
		}

		if strings.ToLower(msg) == "status" {
			// prepare stat struct and send it to the stat accept channel then process it
			stat := StatusOp{respCh: make(chan int)}
			tqs.statusCh <- stat
			task := <-stat.respCh
			if task <= 1 {
				fmt.Fprintf(conn, "%d Task left\n", task)
				continue
			}
			fmt.Fprintf(conn, "%d Tasks left\n", task)
			continue
		}

		if strings.ToLower(msg) == "next" {
			// create a struct with inilized channel to accept the task and return to the client
			next := NextOp{respCh: make(chan *Task)}
			tqs.nextCh <- next
			task := <-next.respCh
			if task == nil {
				fmt.Fprintf(conn, "No task left\n")
				continue
			}
			fmt.Fprintf(conn, "task %d: %s\n", task.ID, task.Content)
			continue
		}

		if strings.HasPrefix(strings.ToLower(msg), "task: ") {
			// parse the task name
			msg = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(msg), "task: "))
			if msg == "" {
				fmt.Fprintf(conn, "Invalid Task\n")
				continue
			}
			// create struct with receiver channel and task to send it to the stateful goroutine
			submit := SubmitOp{task: Task{Content: msg}, respCh: make(chan bool)}
			tqs.submitCh <- submit
			if <-submit.respCh {
				fmt.Fprintf(conn, "Task have been submitted successfully\n")
				continue
			}
		}
		fmt.Fprintf(conn, "Unknown command: %s\n", msg)
	}

}

func main() {
	server := NewTaskQueueServer()
	server.Run() // start the server

	server.client.Wait()
	fmt.Println("Server Shutting down..")
	time.Sleep(500 * time.Millisecond) // simulate server shutting down :)
}
