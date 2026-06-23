package server

import (
	"bufio"
	"fmt"
	"net"
	"go-redis/protocol"
	"go-redis/storage"
)

// StartTCP boots the server with both Memory and PubSub engines
func StartTCP(store *storage.RedisStore, ps *storage.PubSub) {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("[SERVER] Client connected!")
		
		// Spin up a new thread for this client
		go handleConnection(conn, store, ps)
	}
}

func handleConnection(conn net.Conn, store *storage.RedisStore, ps *storage.PubSub) {
	// CLEANUP: Unplug them from the network and the PubSub maps when they leave!
	defer conn.Close()
	defer ps.RemoveClient(conn) 

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Route the command, passing the engines and the connection
		response := protocol.ParseCommand(line, conn, store, ps)
		conn.Write([]byte(response))
	}
	
	fmt.Println("[SERVER] Client disconnected.")
}