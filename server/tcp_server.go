package server

import (
	"bufio"
	"fmt"
	"net"
	"go-redis/protocol"
	"go-redis/storage"
)

// Global instance of our database store
var store = storage.NewRedisStore()

func StartTCPServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("[SERVER] Listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Println("[SERVER] Client connected!")

		// Spawn a new Goroutine for each client to handle concurrency
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[SERVER] Client disconnected.")
			return
		}

		// Pass the incoming string to the protocol parser
		response := protocol.ParseCommand(line, store)

		// Write the result back to the client
		conn.Write([]byte(response))
	}
}