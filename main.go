package main

import (
	"fmt"
	"net"

	"github.com/kjj1998/kvstore/handler"
	"github.com/kjj1998/kvstore/store"
)

var kv = store.NewStore()

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

		go handler.HandleConnection(conn, kv)
	}
}
