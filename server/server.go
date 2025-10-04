package main
import (
    "bufio"
    "fmt"
    "net"
)
var clients []net.Conn
func main() {
    ln, err := net.Listen("tcp", ":9000")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer ln.Close()

    fmt.Println("Server Started on port 9000");
    for{
        conn,err := ln.Accept()
        if err !=nil{
            fmt.Println("Error accepting another client connection:",err)
            continue
        }
        clients = append(clients, conn)
        fmt.Println("New client connected:",conn.RemoteAddr())
        go handleClient(conn)
    }
}
func handleClient(conn net.Conn) {
    defer func() {
        conn.Close()
        for i, c := range clients {
            if c == conn {
                clients = append(clients[:i], clients[i+1:]...)
                break
            }
        }
        fmt.Println("Client disconnected:", conn.RemoteAddr())
    }()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        msg := scanner.Text()
        fmt.Println("Received from", conn.RemoteAddr(), ":", msg)
        broadcast(msg, conn)
    }
}
func broadcast(msg string, sender net.Conn) {
    for _, c := range clients {
        if c != sender {
            fmt.Fprintln(c, msg)
        }
    }
}
