package main

import (
	"fmt"
	"go-redis/server"
	"go-redis/storage"
)

func main() {
	fmt.Println("Initializing Custom Redis Server...")
	
	// Boot the core engines
	store := storage.NewRedisStore()
	pubsub := storage.NewPubSub()

	fmt.Println("[SERVER] Listening on port 6379...")
	
	// Start the TCP listener
	server.StartTCP(store, pubsub)
}