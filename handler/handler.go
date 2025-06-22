package handler

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/kjj1998/kvstore/errors"
	"github.com/kjj1998/kvstore/store"
)

func HandleConnection(conn net.Conn, store *store.Store) {
	defer conn.Close()

handleLoop:
	for {
		reader := bufio.NewReader(conn)
		line, err := reader.ReadString('\n')
		errors.LogError(err, "Error reading from connection: ")

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
				value := store.Get(commands[1])
				fmt.Fprint(conn, value, "\n")
			}
		case "SET":
			if len(commands) < 3 {
				fmt.Fprint(conn, "SET command requires a key and a value\n")
			} else if len(commands) == 5 {
				timeToLive, err := strconv.Atoi(commands[4])

				if commands[3] != "EX" || err != nil {
					fmt.Fprint(conn, "To set values with a time to live expiry, enter command in the format SET <key> <value> EX <integer>\n")
				}

				store.Set(commands[1], commands[2], timeToLive)
			} else {
				store.Set(commands[1], commands[2], 0)
			}
		case "DEL":
			if len(commands) != 2 {
				fmt.Fprint(conn, "DEL command requires a key\n")
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
