// project to practice Go concepts like goroutines, channels, select, timeouts,
// interfaces, custom errors, structs, methods, mutexes, generic types,
// timers, networking, and tickers, worker pool.

// project is about auction that client place their bid on stock values in 15 second open window
// highest bid will win within 15 second window
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// interface for type of stock markets from different resource
type StockFetcher interface {
	FetchPrice(symbol string) (float64, error)
	Name() string
}

// custom error for fetch error
type FetchError struct {
	Source string
	Reason string
}

// should implement Error() method in order to satisfie error
func (fe *FetchError) Error() string {
	return fmt.Sprintf("From %s Failed: %s\n", fe.Source, fe.Reason)
}

// fetch price resource
type MockStock struct {
	prices map[string]float64
	mu     sync.Mutex
}

// implement FetchPrice method to satisfie interface StockFetcher and also fetch price of the stock
func (ms *MockStock) FetchPrice(symbol string) (float64, error) {
	// to avoid any race condition use mutex
	ms.mu.Lock()
	defer ms.mu.Unlock()
	price, ok := ms.prices[symbol]
	if !ok {
		return 0, &FetchError{Source: ms.Name(), Reason: "Can't find symbol"}
	}

	// to visualize change of stock in seconds
	price += float64(time.Now().Nanosecond()%10) / 10
	return price, nil
}

// return which resource we are using to fetch because we are using interface
// so we can't tell which resource we are using that why Name exists
func (ms *MockStock) Name() string {
	return "Mock Stock"
}

// cache each fetch price use mu so to avoid race condition
type PriceCache[T any] struct {
	mu    sync.Mutex
	cache map[string]T
}

// initialize PriceCache
func NewPriceCache[T any]() *PriceCache[T] {
	return &PriceCache[T]{cache: make(map[string]T)}
}

// store price with the symbol we are fetching for example AAPL: 240
func (pc *PriceCache[T]) Store(key string, value T) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.cache[key] = value
}

type Bid struct {
	Symbol string
	Amount float64
}

type BidResult struct {
	Valid  bool
	Amount float64
}

// data structure to manage bid with auction
type Auction struct {
	fetcher     StockFetcher
	currentBid  Bid
	bidCh       chan Bid
	bidResultCh chan BidResult
	cache       *PriceCache[float64]
	mu          sync.Mutex
}

// initilize new Auction
func NewAuction(fetcher StockFetcher) *Auction {
	return &Auction{
		fetcher:     fetcher,
		bidCh:       make(chan Bid, 10),
		bidResultCh: make(chan BidResult, 10),
		cache:       NewPriceCache[float64](),
	}
}

func (a *Auction) BidWorkers(id int, wgWorkers *sync.WaitGroup) {
	defer wgWorkers.Done()
	for bid := range a.bidCh {
		fmt.Printf("Worker %d processing bid $%.2f for %s\n", id, bid.Amount, bid.Symbol)
		time.Sleep(1 * time.Second)
		a.mu.Lock()
		valid := bid.Amount > a.currentBid.Amount
		a.mu.Unlock()
		a.bidResultCh <- BidResult{Valid: valid, Amount: bid.Amount}
	}
}

// method to start the auction
func (a *Auction) Run(wg *sync.WaitGroup, signal chan<- bool) {
	defer wg.Done()
	fmt.Println("Auction started..")
	sym := a.currentBid.Symbol
	fmt.Printf("Cunrrent bid $%.2f of %s\n", a.currentBid.Amount, sym)
	a.cache.Store(sym, a.currentBid.Amount)

	// we use ticker here to visualize the change of price in stock
	ticker := time.NewTicker(2 * time.Second)
	// close ticker too so it will not lick resources
	defer ticker.Stop()
	timer := time.NewTimer(15 * time.Second) // 15 second window for bid

	// worker pools that process each bid
	var wgWorkers sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wgWorkers.Add(1)
		go a.BidWorkers(i, &wgWorkers)
	}

	for {
		// while loop to listen for if new bid comes
		select {
		case bid := <-a.bidResultCh:
			a.mu.Lock()
			if bid.Valid {
				fmt.Printf("New bid $%.2f\n", bid.Amount)
				a.cache.Store(sym, bid.Amount)
				a.currentBid.Amount = bid.Amount
			}
			a.mu.Unlock()
			// if the timer ends at this time so we will drain what channel so it will not take any storage
			if !timer.Stop() {
				// if timer still not end timer.Stop() will return true which mean it able to stop it
				// but if timer ends timer.Stop() will return false which means there is no timer to stop
				// if time stops and return true negation will change to false which means else part execute which
				// mean it will reset
				<-timer.C
			} else {
				timer.Reset(15 * time.Second)
			}
		case <-ticker.C:
			// update the price in every 2 second to visualize stock market

			// i'm forced to use switch because interface type didn't know about
			// the field i am trying to update interface only knows the method it satisfies
			// so switch helps to update the field
			switch v := a.fetcher.(type) {
			case *MockStock:
				v.prices[sym] = a.currentBid.Amount
			default:
				fmt.Println("no match")
			}
			price, err := a.fetcher.FetchPrice(sym)
			if err != nil {
				fmt.Printf("Error fetching %s\n", sym)
				continue
			}

			fmt.Printf("New price for %s is $%.2f\n", sym, price)
		case <-timer.C:
			// if timer is up send true to signal channel to let know other func that the timer is up
			// and close the channel so there will be no data sent to it
			signal <- true
			close(signal)
			return
		}
	}
}

func main() {
	// declare and initilize
	var wg sync.WaitGroup
	fetcher := &MockStock{
		prices: map[string]float64{"AAPL": 230.01, "GOOGL": 150.02},
	}
	auction := NewAuction(fetcher)
	sym := "AAPL"
	signal := make(chan bool, 1)

	price, err := fetcher.FetchPrice(sym)
	if err != nil {
		fmt.Printf("Error Fetching %s - %s\n", fetcher.Name(), err.Error())
		return
	}

	auction.currentBid = Bid{Symbol: sym, Amount: price}

	wg.Add(1)
	go auction.Run(&wg, signal)

	// create the server that the client listen to and send the bid
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// reason to create this channels is the behaviour ln.Accept() because this will block the main
	// goroutine so to separte the concern i used a channel to accept net.Conn object and run ln.Accept()
	// in goroutine so it will not block the main goroutine when ln.Accept() accepts or finds a client it will send it
	// through the channel
	connCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				errCh <- err
				return
			}
			connCh <- conn
		}
	}()

	// use this flag to separate concern from getting client or no client there
	var flag bool = false

	for {
		select {
		case conn := <-connCh:
			flag = true
			wg.Add(1)
			// start the goroutine HandleConnection once we have a client to connet with
			go HandleConnection(conn, &wg, auction, signal)
		case err := <-errCh:
			fmt.Println("Connection error: ", err)
			close(signal)
			return
		case <-signal:
			if !flag {
				fmt.Println("No client...")
				fmt.Printf("Auction Ended.... winner $%.2f for %s\n", auction.currentBid.Amount, auction.currentBid.Symbol)
			}
			return
		}
	}
	wg.Wait()
}

func HandleConnection(conn net.Conn, wg *sync.WaitGroup, auction *Auction, signal <-chan bool) {
	defer wg.Done()
	defer conn.Close()
	// Read anything that comes from the connection
	reader := bufio.NewReader(conn)
	// use this channel to receive the amount of bid from the client
	inputCh := make(chan string, 1)

	// reason i used goroutine here is that because of reader.ReadString behaviour
	// it will block the routine until received value from the client means that
	// even if the window is closed after 15 second the code is hang here waiting to receive value
	// solution: use goroutine to handle using concurrency
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				close(inputCh)
				return
			}
			inputCh <- strings.TrimSpace(input)
		}
	}()

	for {
		select {
		case <-signal:
			// if the times up we will receive signal from signal channel
			fmt.Printf("Bidding Ended.... winner $%.2f for %s\n", auction.currentBid.Amount, auction.currentBid.Symbol)
			fmt.Fprintln(conn, "EXIT")
			return
		case input, ok := <-inputCh:
			// receive bid amount or other signal from the client by this channel
			if !ok {
				fmt.Println("Input channel closed")
				os.Exit(1)
				return
			}
			if input == "exit" {
				return
			}
			// typecast the value receive from the channel
			bid, err := strconv.Atoi(input)
			if err != nil {
				fmt.Fprintf(conn, "Invalid bid: %v\n", err)
				continue
			}
			// here the reason used goroutine is that if the auction.bidCh channel is already full it will block the
			// routine and also if the auction system is busy processing other bids because bids comes from different clients
			// so solution to use goroutine to handle it from another goroutine and continue the process
			wg.Add(1)
			go func() {
				defer wg.Done()
				auction.bidCh <- Bid{Symbol: "AAPL", Amount: float64(bid)}
			}()
		}
	}
}
