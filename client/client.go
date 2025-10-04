package main
import (
    "bufio"
    "fmt"
    "net"
    "os"
)
func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:9000")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()
    fmt.Println("Connected to server. Type messages:")
    go func() {
        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            fmt.Println("Server:", scanner.Text())
        }
        fmt.Println("Disconnected from server.")
        os.Exit(0)
    }()
    stdin := bufio.NewScanner(os.Stdin)
    for stdin.Scan() {
        fmt.Fprintln(conn, stdin.Text())
    }
}
