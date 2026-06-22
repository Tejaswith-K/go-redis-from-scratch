package protocol

import (
	"strings"
	"go-redis/storage"
)

// ParseCommand cleans up network whitespace and routes the command
func ParseCommand(rawLine string, store *storage.RedisStore) string {
	cleanLine := strings.TrimSpace(rawLine)
	parts := strings.Fields(cleanLine)

	if len(parts) == 0 {
		return "-ERR empty command\r\n"
	}

	command := strings.ToUpper(parts[0])

	switch command {
	case "SET":
		if len(parts) < 3 {
			return "-ERR wrong number of arguments for 'set' command\r\n"
		}
		key := parts[1]
		value := strings.Join(parts[2:], " ")
		store.Set(key, value)
		return "+OK\r\n"

	case "GET":
		if len(parts) != 2 {
			return "-ERR wrong number of arguments for 'get' command\r\n"
		}
		key := parts[1]
		value, exists := store.Get(key)
		if !exists {
			return "$-1\r\n" // Standard Redis response for Not Found (Nil)
		}
		return "+" + value + "\r\n"

	case "PING":
		return "+PONG\r\n"

	default:
		return "-ERR unknown command '" + command + "'\r\n"
	}
}