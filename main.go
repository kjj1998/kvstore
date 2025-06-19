package main

import (
	"fmt"
	"net"
)

var store = NewStore()

func main() {
	ln, err := net.Listen("tcp", ":8080")
	fmt.Println("Server started and listening on port 8080...")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, store)
	}
}
