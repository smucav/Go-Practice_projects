# Go TCP Broadcast Server

This is a simple TCP-based broadcast server and client implemented in Go. The server listens for incoming client connections, receives messages from clients, and broadcasts them to all connected clients.

## Features
- Supports multiple concurrent clients.
- Uses **goroutines** for handling multiple clients simultaneously.
- Implements **rate limiting** to prevent spam.
- Clients can send messages, and all connected clients receive them in real time.
- Clients can type `exit` to disconnect.

## Installation & Usage

### Prerequisites
- Go (1.18 or later)

### Clone the Repository
```sh
git clone https://github.com/your-username/go-broadcast-server.git
cd go-broadcast-server
```

### Run the Server
```sh
go run broadcast_server.go
```

### Run a Client
In a new terminal window, execute:
```sh
go run broadcast_client.go
```
You can run multiple clients in separate terminal windows to test broadcasting.

### Sending Messages
Once connected, type a message and press **Enter** to send it. All connected clients will receive the message.

To disconnect, type:
```sh
exit
```

## Project Structure
```
├── broadcast_server.go  # TCP Broadcast Server
├── broadcast_client.go  # TCP Client
├── README.md            # Documentation
```

## How It Works
1. The server listens on `localhost:8080` for incoming connections.
2. When a client connects, the server assigns it a unique ID and listens for messages.
3. Messages from clients are broadcasted to all connected clients.
4. The server implements **rate limiting** using a bursty limiter to prevent spam.
5. Clients can disconnect by sending `exit`.

## Example Usage
**Start the Server:**
```
$ go run broadcast_server.go
Server started on :8080
```

**Connect Clients:**
```
$ go run broadcast_client.go
$ Hello from Client 1
```

**Client 2 Output:**
```
Client 1: Hello from Client 1
```

## Improvements & Next Steps
- Add **authentication** for secure client connections.
- Implement **message history**.
- Use **WebSockets** instead of raw TCP for a web-friendly solution.

## License
This project is open-source and available under the MIT License.

---
Enjoy building with Go! 🚀
