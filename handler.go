package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
		}

		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Client requested to exit.")
			break
		}

		fmt.Printf("Received: %s", line)
		fmt.Fprint(conn, "Message received\n")
	}
}
