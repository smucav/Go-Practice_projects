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

type Score struct {
	GameId string    `json:"game_id"`
	TeamA  string    `json:"team_a"`
	TeamB  string    `json:"team_b"`
	ScoreA int       `json:"score_a"`
	ScoreB int       `json:"score_b"`
	Update time.Time `json:"updated"`
}

type ScoreReport struct {
	data   string
	respCh chan string
}

type ScoreBoard struct {
	reportCh    chan ScoreReport
	scores      map[string]Score
	mu          sync.Mutex
	broadcastCh chan Score
	clients     sync.WaitGroup
	shutdownCh  chan struct{}
}

func NewScoreBoard() *ScoreBoard {
	return &ScoreBoard{
		scores:      make(map[string]Score),
		broadcastCh: make(chan Score, 10),
		reportCh:    make(chan ScoreReport),
		shutdownCh:  make(chan struct{}),
	}
}

func (sb *ScoreBoard) Run() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Server Error")
		return
	}

	defer ln.Close()
	fmt.Println("ScoreBoard server started on :8080")

	go func() {
		for {
			select {
			case report := <-sb.reportCh:
				var score Score
				err := json.Unmarshal([]byte(report.data), &score)
				if err != nil {
					report.respCh <- fmt.Sprintf("Invalid score data: %v\n", err)
					continue
				}
				if score.ScoreA < 0 || score.ScoreB < 0 {
					panic("Negative scores not allowed")
				}
				score.Update = time.Now()
				sb.mu.Lock()
				sb.scores[score.GameId] = score
				sb.mu.Unlock()
				sb.broadcastCh <- score
				report.respCh <- "Report Submitted Successfully\n"
			case <-sb.shutdownCh:
				fmt.Println("Score aggregator shutting down..")
				return
			}
		}
	}()

	go func() {
		timer := time.NewTicker(10 * time.Second)
		defer timer.Stop()
		for {
			select {
			case score := <-sb.broadcastCh:
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

var activeViewers = make(map[net.Conn]struct{})
var viewMu sync.Mutex

func (sb *ScoreBoard) broadcast(score Score) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	for conn := range activeViewers {
		fmt.Fprintf(conn, "Game: %s |%s %d - %d %s| %v\n", score.GameId, score.TeamA, score.ScoreA, score.ScoreB, score.TeamB, score.Update)
	}
}

func (sb *ScoreBoard) handleClient(conn net.Conn) {
	defer sb.clients.Done()
	defer conn.Close()

	fmt.Fprintf(conn, "welcome to score board use Command like 'VIEW' 'to report use 'REPORT: <report>\n")

	reader := bufio.NewReader(conn)

	for {
		defer func() {
			if f := recover(); f != nil {
				fmt.Fprintf(conn, "Panic caught: %v\n", f)
				fmt.Println("Client %s panicked but recovered: %v\n", conn.RemoteAddr().String(), f)
			}
		}()

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client Error sending")
			viewMu.Lock()
			delete(activeViewers, conn)
			viewMu.Unlock()
			return
		}

		msg := strings.TrimSpace(input)

		if msg == "EXIT" {
			fmt.Fprintf(conn, "Good bye :)\n")
			viewMu.Lock()
			delete(activeViewers, conn)
			viewMu.Unlock()
			fmt.Printf("client %s exited\n", conn.RemoteAddr().String())
			return
		}
		if msg == "VIEW" {
			fmt.Fprintf(conn, "Now you are View scores\n")
			viewMu.Lock()
			activeViewers[conn] = struct{}{}
			viewMu.Unlock()
			fmt.Fprintf(conn, "Now viewing live socres\n")
			continue
		}

		if strings.HasPrefix(msg, "REPORT ") {
			report := strings.TrimSpace(strings.TrimPrefix(msg, "REPORT "))

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
