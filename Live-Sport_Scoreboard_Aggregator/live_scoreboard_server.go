// Project to learn goroutine, channel, panic and recover mutex struct, net.. more
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// board where each game information stored
type Score struct {
	GameId string    `json:"game_id"`
	TeamA  string    `json:"team_a"`
	TeamB  string    `json:"team_b"`
	ScoreA int       `json:"score_a"`
	ScoreB int       `json:"score_b"`
	Update time.Time `json:"updated"` // used to show to which time is the score showing
}

// to accept game information from reporter and send acceptance to the reporter
type ScoreReport struct {
	data   string // game information
	respCh chan string
}

// Score board which have games information and channel to broadcast game information to viewers and also keep track of clients activity by goroutine
type ScoreBoard struct {
	reportCh    chan ScoreReport
	scores      map[string]Score
	mu          sync.Mutex
	broadcastCh chan Score
	clients     sync.WaitGroup
	shutdownCh  chan struct{}
}

// initialize the scoreboard
func NewScoreBoard() *ScoreBoard {
	return &ScoreBoard{
		scores:      make(map[string]Score),
		broadcastCh: make(chan Score, 10),
		reportCh:    make(chan ScoreReport),
		shutdownCh:  make(chan struct{}),
	}
}

// run the server
func (sb *ScoreBoard) Run() {
	// listen to any connection
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Server Error")
		return
	}

	defer ln.Close()
	fmt.Println("ScoreBoard server started on :8080")

	// acceptor goroutine from reporters
	go func() {
		for {
			select {
			case report := <-sb.reportCh:
				// used immediate invoked function expression(IIFE)
				// because of panic and recover per report not to affect other reports
				func() {
					defer func() {
						if f := recover(); f != nil {
							report.respCh <- fmt.Sprintf("panic negative score not allowed\n")
						}
					}()

					var score Score
					err := json.Unmarshal([]byte(report.data), &score)
					if err != nil {
						report.respCh <- fmt.Sprintf("Invalid score data: %v\n", err)
						return
					}
					if score.ScoreA < 0 || score.ScoreB < 0 {
						panic("Negative scores not allowed")
					}
					score.Update = time.Now() // set time to which time the score is begin updated
					// use lock and unlock to avoid race condition
					sb.mu.Lock()
					// store it so we can send the score again to viewers when no new updated score is not begin received from reporter
					sb.scores[score.GameId] = score
					sb.mu.Unlock()
					// send the score that is begin received from report to broadcast goroutine to send to all viewers
					sb.broadcastCh <- score
					report.respCh <- "Report Submitted Successfully\n"
				}()
			case <-sb.shutdownCh:
				fmt.Println("Score aggregator shutting down..")
				return
			}
		}
	}()

	// broadcastor
	go func() {
		// use ticker here to send each score again when no new updated score result is not begin received
		timer := time.NewTicker(10 * time.Second)
		defer timer.Stop()
		for {
			select {
			case score := <-sb.broadcastCh:
				// send to the function that will broadcast to active viewers
				sb.broadcast(score)
			case <-timer.C:
				sb.mu.Lock()
				for _, score := range sb.scores {
					sb.broadcast(score)
				}
				sb.mu.Unlock()
			case <-sb.shutdownCh:
				fmt.Println("Broadcast shutting down..")
				return
			}
		}
	}()

	// handle each new client
	sb.clients.Add(1)
	go func() {
		defer sb.clients.Done()
		for {
			select {
			case <-sb.shutdownCh:
				return
			default:
				conn, err := ln.Accept()
				if err != nil {
					fmt.Println("Server down")
					return
				}
				sb.clients.Add(1)
				go sb.handleClient(conn)
			}
		}

	}()

	sb.clients.Wait()
	fmt.Println("before Closing")
	close(sb.shutdownCh)
}

// where all active viewers stored
// once they exited they will be deleted from it
var activeViewers = make(map[net.Conn]struct{})
var viewMu sync.Mutex

// broadcast to each viewers
func (sb *ScoreBoard) broadcast(score Score) {
	viewMu.Lock()
	defer viewMu.Unlock()
	for conn := range activeViewers {
		fmt.Fprintf(conn, "Game: %s |%s %d - %d %s| %v\n", score.GameId, score.TeamA, score.ScoreA, score.ScoreB, score.TeamB, score.Update)
	}
}

func (sb *ScoreBoard) handleClient(conn net.Conn) {
	defer sb.clients.Done()
	defer conn.Close()

	fmt.Fprintf(conn, "welcome to score board use Command like 'VIEW' 'to report use 'REPORT: <report>\n")

	// reader from the connection
	reader := bufio.NewReader(conn)

	for {

		// use this recover to recover from panic happend when receiving bad report
		defer func() {
			if f := recover(); f != nil {
				fmt.Fprintf(conn, "Panic caught: %v\n", f)
				fmt.Println("Client %s panicked but recovered: %v\n", conn.RemoteAddr().String(), f)
			}
		}()

		// read client input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client Error sending")
			viewMu.Lock()
			delete(activeViewers, conn)
			viewMu.Unlock()
			return
		}

		// remove trailing spaces
		msg := strings.TrimSpace(input)

		if msg == "EXIT" {
			// handle exit and also delete from active viewers
			fmt.Fprintf(conn, "Good bye :)\n")
			viewMu.Lock()
			delete(activeViewers, conn)
			viewMu.Unlock()
			fmt.Printf("client %s exited\n", conn.RemoteAddr().String())
			return
		}

		if msg == "VIEW" {
			// activate the client as viewer
			viewMu.Lock()
			activeViewers[conn] = struct{}{}
			viewMu.Unlock()
			fmt.Fprintf(conn, "Now viewing live socres\n")
			continue
		}

		if strings.HasPrefix(msg, "REPORT ") {
			// accept a report
			report := strings.TrimSpace(strings.TrimPrefix(msg, "REPORT "))
			if report == "{bad}" {
				panic("invalid report")
			}
			// ship it to send to the acceptor to process and broadcast it
			reportsend := ScoreReport{
				data:   report,
				respCh: make(chan string),
			}

			sb.reportCh <- reportsend
			fmt.Fprintf(conn, <-reportsend.respCh)
			continue
		}
		fmt.Fprintf(conn, "Unknown command: %s\n", msg)
	}
}

func main() {
	server := NewScoreBoard()

	server.Run()
}
