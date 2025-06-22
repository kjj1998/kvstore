package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kjj1998/kvstore/errors"
	"github.com/kjj1998/kvstore/handler"
	"github.com/kjj1998/kvstore/store"
)

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	var kv = store.NewStore(context, cancel)

	ln, err := net.Listen("tcp", ":8080")
	fmt.Println("Server started and listening on port 8080...")
	errors.LogError(err, "Error starting server: ")

	kv.RecoverFromLog()
	os.Truncate("persistent_log.txt", 0)

	kv.StartWALWriterGoroutine()
	kv.BackgroundCleanUpService(30 * time.Second)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down gracefully...")
		cancel()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if os.IsTimeout(err) || strings.Contains(err.Error(), "use of closed network connection") {
				break // Expected during shutdown — do not log
			}
			errors.LogError(err, "Error accepting connection: ")
			continue
		}
		go handler.HandleConnection(conn, kv)
	}

	fmt.Println("Server shutdown complete")
}
