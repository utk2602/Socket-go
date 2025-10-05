package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mu      sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("🚀 Minefield Escape v0.0.1 server started on port 8080")

	game := NewGame(5) // create 5x5 grid

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		playerName := fmt.Sprintf("Player%d", len(clients)+1)
		clients[conn] = playerName
		mu.Unlock()

		fmt.Printf("✅ %s connected.\n", playerName)
		go handleConnection(conn, game, playerName)
	}
}

func handleConnection(conn net.Conn, game *Game, playerName string) {
	defer func() {
		conn.Close()
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		fmt.Printf("❌ %s disconnected.\n", playerName)
	}()

	conn.Write([]byte("Welcome to Minefield Escape!\n"))
	conn.Write([]byte("Commands: MOVE UP/DOWN/LEFT/RIGHT or SHOW\n"))

	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.ToUpper(message) == "SHOW" {
			conn.Write([]byte(game.Display(playerName)))
			continue
		}

		if strings.HasPrefix(strings.ToUpper(message), "MOVE") {
			direction := strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(message), "MOVE"))
			resp := game.MovePlayer(playerName, direction)
			broadcast(resp + "\n")
			continue
		}

		conn.Write([]byte("Unknown command\n"))
	}
}

func broadcast(msg string) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		conn.Write([]byte(msg))
	}
}
