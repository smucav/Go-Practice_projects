# Task Queue Server

## Overview
This project implements a simple Task Queue Server using Golang. It allows clients to submit tasks, retrieve the next task, and check the queue status via a TCP connection. The server efficiently manages task submissions and retrievals using Go channels and goroutines, demonstrating concurrency and networking principles.

## Features
- **Submit tasks**: Clients can add tasks to the queue.
- **Retrieve tasks**: Clients can request the next task in FIFO order.
- **Check queue status**: Clients can check how many tasks are left in the queue.
- **Concurrent client handling**: Multiple clients can interact with the server simultaneously.
- **Graceful shutdown**: Server handles client disconnections properly.

## Technologies Used
- **Golang**
- **TCP Networking**
- **Goroutines and Channels**
- **Synchronization with WaitGroups**

## Installation
### Prerequisites
- Install [Go](https://golang.org/doc/install)

### Steps to Run
1. Clone this repository:
   ```sh
   git clone https://github.com/your-repo/task-queue-server.git
   cd task-queue-server
   ```
2. Build and start the server:
   ```sh
   go run task_queue_server.go
   ```
3. In another terminal, start a client:
   ```sh
   go run task_queue_client.go
   ```
4. Follow the client instructions to interact with the server.

## Usage
### Commands
- `Task: <task>` - Submits a new task.
- `NEXT` - Retrieves the next available task.
- `STATUS` - Checks how many tasks are left in the queue.
- `EXIT` - Disconnects the client from the server.

## Topics Covered
This project is a great learning exercise for practicing:
- **Golang Concurrency**: Using goroutines and channels to manage a task queue.
- **TCP Server-Client Communication**: Handling connections using the `net` package.
- **Synchronization Mechanisms**: Using `sync.WaitGroup` to manage multiple clients.
- **Stateful Goroutines**: Keeping state within a goroutine to process incoming tasks.
- **Error Handling in Network Programming**: Managing client disconnections and errors.

## Future Enhancements
- Add persistent task storage using a database.
- Implement priority-based task retrieval.
- Create a web interface for easier task management.
- Improve error handling and logging.

## License
This project is for learning purposes and is open-source. Feel free to modify and improve it!

## Author
Developed by **Daniel Tujuma**

