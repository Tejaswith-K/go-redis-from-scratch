package main

import (
	"fmt"
	"go-redis/server"
)

func main() {
	fmt.Println("Initializing Custom Redis Server...")
	
	// Start the TCP server on standard Redis port 6379
	err := server.StartTCPServer("6379")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}