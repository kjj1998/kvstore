package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/kjj1998/kvstore/store"
)

func HandleConnection(conn net.Conn, store *store.Store) {
	defer conn.Close()

handleLoop:
	for {
		reader := bufio.NewReader(conn)

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
		}

		commands := strings.Fields(line)
		if len(commands) == 0 {
			fmt.Println("No commands received.")
			continue
		}

		switch commands[0] {
		case "GET":
			if len(commands) != 2 {
				fmt.Fprint(conn, "GET command requires a key\n")
			} else {
				if value, exists := store.Get(commands[1]); exists {
					fmt.Fprint(conn, value, "\n")
				} else {
					fmt.Fprint(conn, "NULL\n")
				}
			}
		case "SET":
			if len(commands) != 3 {
				fmt.Fprint(conn, "SET command requires a key and a value\n")
			} else {
				store.Set(commands[1], commands[2])
			}
		case "DELETE":
			if len(commands) != 2 {
				fmt.Fprint(conn, "DELETE command requires a key\n")
			} else {
				store.Delete(commands[1])
			}
		case "EXIT":
			break handleLoop
		default:
			fmt.Fprint(conn, "Unknown command\n")
			fmt.Println("Unknown command received:", commands[0])
		}
	}
}
