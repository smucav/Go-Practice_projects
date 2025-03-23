# Hit Counter Server

A simple TCP-based Hit Counter Server implemented in Go that tracks page visits and allows clients to query hit statistics.

## Features
- Record hits for different pages via TCP commands.
- Retrieve statistics on page hits.
- Supports multiple concurrent clients.

## Requirements
- Go 1.18+

## Installation
Clone the repository and navigate to the project directory:
```sh
$ git clone <your-repo-url>
$ cd hit-counter-server
```

## Usage
### Run the Server
Start the hit counter server:
```sh
$ go run hit_counter_server.go
```
The server will start listening on port `8080`.

### Run the Client
Open another terminal and run:
```sh
$ go run hit_counter_client.go
```
This connects to the server.

### Commands
- `GET <page>` - Record a hit for a page.
- `STATS` - Get the total hits for all tracked pages.
- `exit` - Disconnect from the server.

### Example Usage
1. **Start the server**
```sh
$ go run hit_counter_server.go
```
2. **Start a client**
```sh
$ go run hit_counter_client.go
```
3. **Send commands in the client terminal**
```sh
GET homepage
GET about
STATS
exit
```
4. **Example Server Output**
```sh
Hit Counter Server started on :8080
Recorded hit for homepage from 127.0.0.1:54321 (total: 1)
Recorded hit for about from 127.0.0.1:54321 (total: 1)
```

## Notes
- The server uses `atomic.Uint64` for safe concurrent hit counting.
- Uses `sync.Mutex` to protect access to the page map.

## License
This project is licensed under the MIT License.

