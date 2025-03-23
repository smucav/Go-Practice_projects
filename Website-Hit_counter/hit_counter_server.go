// project to pratice topics like: goroutines, channels, select, timeouts,
// interfaces, custom errors, structs, methods, mutexes, generic types
// timers, networking, tickers, worker pools, and rate limiting.
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type PageHit struct {
	Name string
	// atomic counter with thread-safety without requiring exlicit locking without causing race conditions
	hit atomic.Uint64
}

// struct that keeps the data of pages and their hit
type HitCounterServer struct {
	pages map[string]*PageHit
	mu    sync.Mutex     // this is used for the pages map to avoid race condition not PageHit hit which is atomic counter
	wg    sync.WaitGroup // to keep track of goroutines
}

func NewHitCounterServer() *HitCounterServer {
	return &HitCounterServer{
		pages: make(map[string]*PageHit),
	}
}

func (hcs *HitCounterServer) Run() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Server Error:", err)
		return
	}

	defer ln.Close()
	fmt.Printf("Hit Counter Server started on :8080\n")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Server Disconnected: ", err)
			return
		}
		hcs.wg.Add(1)
		go hcs.handleClient(conn)
	}
}

func (hcs *HitCounterServer) handleClient(conn net.Conn) {
	defer hcs.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn) // connect to the client to read their message
	fmt.Fprintf(conn, "Send 'GET <page>' to record a hit, 'STATS' for totals, or 'exit' to quit:\n")

	for {
		input, err := reader.ReadString('\n') // read client message
		if err != nil {
			fmt.Printf("Client %s Disconnected", conn.RemoteAddr().String())
			return
		}

		cmd := strings.TrimSpace(input)
		if cmd == "exit" {
			fmt.Printf("Client %s exited\n", conn.RemoteAddr().String())
			return
		}

		if cmd == "STATS" {
			hcs.mu.Lock()
			for _, page := range hcs.pages {
				fmt.Fprintf(conn, "Site %s Hit %d\n", page.Name, page.hit.Load())
			}
			hcs.mu.Unlock()
			continue
		}

		// check if client wants to visite website
		if strings.HasPrefix(cmd, "GET ") {
			ph := strings.TrimSpace(strings.TrimPrefix(cmd, "GET ")) // parse the website name
			if ph == "" {
				fmt.Fprintf(conn, "Invalid Page\n")
				continue
			}

			page := hcs.getOrCreatePage(ph) // check if page exists and return page data
			page.hit.Add(1)
			fmt.Printf("Recorded hit for %s from %s (total: %d)\n", page.Name, conn.RemoteAddr().String(), page.hit.Load())
			fmt.Fprintf(conn, "Hit recorded for %s\n", page.Name)
			continue
		}
		fmt.Fprintf(conn, "Unknown command: %s\n", cmd)
	}
}

func (hcs *HitCounterServer) getOrCreatePage(page string) *PageHit {
	hcs.mu.Lock()
	defer hcs.mu.Unlock()

	if pg, exists := hcs.pages[page]; exists {
		return pg // return existed page if their
	}

	pg := &PageHit{Name: page} // create new page and return it's data
	hcs.pages[page] = pg
	return pg
}

func main() {
	server := NewHitCounterServer()
	server.Run()
}
