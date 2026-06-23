package protocol

import (
	"net"
	"strconv"
	"strings"
	"go-redis/storage"
)

// ParseCommand now takes the live net.Conn and the PubSub engine
func ParseCommand(rawLine string, conn net.Conn, store *storage.RedisStore, ps *storage.PubSub) string {
	cleanLine := strings.TrimSpace(rawLine)
	parts := strings.Fields(cleanLine)

	if len(parts) == 0 {
		return "-ERR empty command\r\n"
	}

	command := strings.ToUpper(parts[0])

	switch command {
	// --- CORE DATABASE ---
	case "SET":
		if len(parts) < 3 {
			return "-ERR wrong number of arguments for 'set'\r\n"
		}
		store.Set(parts[1], strings.Join(parts[2:], " "))
		return "+OK\r\n"

	case "GET":
		if len(parts) != 2 {
			return "-ERR wrong number of arguments for 'get'\r\n"
		}
		value, exists := store.Get(parts[1])
		if !exists {
			return "$-1\r\n"
		}
		return "+" + value + "\r\n"

	case "PING":
		return "+PONG\r\n"

	// --- TTL / EXPIRATION ---
	case "EXPIRE":
		if len(parts) < 3 {
			return "-ERR wrong number of arguments for 'expire'\r\n"
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return "-ERR value is not an integer\r\n"
		}
		if store.Expire(parts[1], seconds) {
			return ":1\r\n"
		}
		return ":0\r\n"

	case "TTL":
		if len(parts) != 2 {
			return "-ERR wrong number of arguments for 'ttl'\r\n"
		}
		return ":" + strconv.FormatInt(store.TTL(parts[1]), 10) + "\r\n"

	// --- PUB/SUB MESSAGING ---
	case "SUBSCRIBE":
		if len(parts) != 2 {
			return "-ERR wrong number of arguments for 'subscribe'\r\n"
		}
		channel := parts[1]
		ps.Subscribe(channel, conn) // Save their socket!
		return "+Subscribed to " + channel + "\r\n"

	case "PUBLISH":
		if len(parts) < 3 {
			return "-ERR wrong number of arguments for 'publish'\r\n"
		}
		channel := parts[1]
		message := strings.Join(parts[2:], " ")
		
		receivers := ps.Publish(channel, message)
		
		// Return integer: number of clients that received it
		return ":" + strconv.Itoa(receivers) + "\r\n"

	default:
		return "-ERR unknown command '" + command + "'\r\n"
	}
}