package game

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func game(){
	conn,err := net.Dial("tcp","localhost:9000")
	if err!=nil{
		fmt.Println("Error connecting to the server:",err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to the server. Type messages:")
}
