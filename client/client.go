package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := consoleReader.ReadString('\n')
		conn.Write([]byte(text))
	}
}
