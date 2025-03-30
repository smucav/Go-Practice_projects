# Live Scoreboard Server

## Overview
The **Live Scoreboard Server** is a Go-based project designed to help developers learn and practice key Go concepts, such as goroutines, channels, mutexes, panic recovery, and network programming. This project implements a real-time scoreboard system that allows reporters to submit game scores and viewers to receive live updates.

## Key Concepts Covered
- **Goroutines**: Used to handle multiple clients and background tasks concurrently.
- **Channels**: Facilitates communication between different components of the system.
- **Mutexes**: Ensures safe concurrent access to shared resources.
- **Panic and Recover**: Prevents the entire system from crashing due to bad inputs.
- **TCP Networking**: Implements a simple client-server model using Go's `net` package.

## How It Works
The project consists of two main components:
1. **Live Scoreboard Server (`live_scoreboard_server.go`)**
   - Accepts game score reports from reporters.
   - Broadcasts updated scores to active viewers.
   - Uses a ticker to periodically resend the latest scores to viewers.
   - Handles client disconnections and errors gracefully.

2. **Client (`live_scoreboard_client.go`)**
   - Connects to the server as a viewer or reporter.
   - Can send commands to view live scores or submit new scores.
   - Handles server responses asynchronously.

## Features
- **Real-time updates**: Reported scores are instantly broadcasted to all viewers.
- **Graceful error handling**: Uses `panic` and `recover` to catch and handle runtime errors.
- **Concurrency**: Uses goroutines to manage multiple clients simultaneously.
- **Thread-safe operations**: Employs mutex locking to prevent data race conditions.
- **Custom client commands**:
  - `VIEW`: Subscribes the client to live score updates.
  - `REPORT <game_data>`: Submits a new game score to the server.
  - `EXIT`: Disconnects from the server.

## Installation & Usage
### 1. Clone the Repository
```sh
 git clone https://github.com/yourusername/Go-Practice_Projects.git
 cd Go-Practice_Projects
```

### 2. Start the Scoreboard Server
```sh
 go run live_scoreboard_server.go
```

### 3. Run a Client
```sh
 go run live_scoreboard_client.go
```

### 4. Example Commands
#### As a Reporter:
```sh
 REPORT {"game_id":"001","team_a":"TeamX","team_b":"TeamY","score_a":2,"score_b":3}
```
#### As a Viewer:
```sh
 VIEW
```

## Contribution
This project is part of the **Go Practice Projects** repository, aimed at learning and mastering Go concepts. Contributions are welcome! Feel free to:
- Improve the project structure.
- Optimize concurrency handling.
- Add more features like persistent storage or a web interface.

## License
This project is open-source and available under the MIT License.

## Author
Developed by **[smucav](https://github.com/smucav)** and contributors. Let's build and learn Go together!

