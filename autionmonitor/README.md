# Auction Bid Monitor

## Overview
Auction Bid Monitor is a Go-based auction system that allows clients to place bids on stock values within a 15-second open window. The highest bid within this timeframe wins. This project demonstrates various Go concepts, including:

- Goroutines & Channels
- `select` & Timeouts
- Interfaces & Custom Errors
- Structs & Methods
- Mutexes & Concurrency
- Generic Types
- Timers & Tickers
- Worker Pool
- Networking

## Features
- Mock stock price fetching with real-time updates.
- Concurrent bid processing using worker pools.
- TCP server-client communication for bid placement.
- Auction auto-ends after 15 seconds.
- Cache implementation for stock prices.

## Installation & Setup
### Prerequisites
- Go 1.18+

### Clone the Repository
```sh
git clone https://github.com/yourusername/auction-bid-monitor.git
cd auction-bid-monitor
```

### Running the Server
```sh
go run auctionbidmonitor.go
```

### Running the Client
Open another terminal and run:
```sh
go run auction_client.go
```

You can also run multiple clients to simulate multiple bidders.

## Usage
1. Start the server (`auctionbidmonitor.go`).
2. Start a client (`auction_client.go`).
3. Enter bid amounts in the client terminal.
4. The auction ends after 15 seconds, and the highest bid wins.

## Example Interaction
```
$ go run auctionbidmonitor.go
Auction started..
Current bid $230.01 of AAPL
New price for AAPL is $230.10
...
Auction Ended.... winner $250.00 for AAPL
```

### Client Example
```
$ go run auction_client.go
Enter bid (or type 'exit' to quit): 240
Enter bid (or type 'exit' to quit): 250
Server: New bid $250.00
...
Server: Bidding Ended.... winner $250.00 for AAPL
Server is shutting down... Exiting client.
```

## Code Structure
- `auctionbidmonitor.go`: The main auction server logic.
- `auction_client.go`: A TCP client to place bids.

## Improvements & Future Enhancements
- Implement real-time stock data fetching from external APIs.
- Add support for multiple stocks in a single auction.
- Improve bid validation and error handling.
- Introduce a web UI for bidding.

## License
This project is licensed under the MIT License.

---

Happy coding! 🚀

