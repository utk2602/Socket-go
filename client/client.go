package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var username = "Anonymous"

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printHelp() {
	fmt.Println(`
📜 Available commands:
----------------------------------
/help                → Show this help menu
/rename <name>       → Change your display name
/whisper <name> <msg> → Send a private message
/clear               → Clear your screen
/exit                → Disconnect from server
----------------------------------
`)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("❌ Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ Connected to server!")
	fmt.Println("💬 Type '/help' to see all available commands.\n")

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("\n📨 " + scanner.Text())
			fmt.Print("> ")
		}
		fmt.Println("\n⚠️ Server connection closed.")
		os.Exit(0)
	}()

	consoleReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, err := consoleReader.ReadString('\n')
		if err != nil {
			fmt.Println("⚠️ Error reading input:", err)
			break
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		switch {
		case strings.HasPrefix(text, "/exit"):
			fmt.Println("👋 Disconnecting...")
			return

		case strings.HasPrefix(text, "/clear"):
			clearScreen()

		case strings.HasPrefix(text, "/help"):
			printHelp()

		case strings.HasPrefix(text, "/rename "):
			newName := strings.TrimSpace(strings.TrimPrefix(text, "/rename "))
			if newName == "" {
				fmt.Println("⚠️ Usage: /rename <name>")
				continue
			}
			username = newName
			fmt.Printf("✅ You are now known as '%s'\n", username)

		case strings.HasPrefix(text, "/whisper "):
			msg := strings.TrimSpace(strings.TrimPrefix(text, "/whisper "))
			if msg == "" {
				fmt.Println("⚠️ Usage: /whisper <name> <message>")
				continue
			}
			conn.Write([]byte(fmt.Sprintf("/whisper %s\n", msg)))

		default:
			conn.Write([]byte(fmt.Sprintf("[%s]: %s\n", username, text)))
		}
	}
}
