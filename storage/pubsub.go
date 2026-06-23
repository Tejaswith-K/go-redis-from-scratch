package storage

import (
	"net"
	"sync"
)

// PubSub tracks live connections listening to specific channels
type PubSub struct {
	mu       sync.RWMutex
	channels map[string][]net.Conn
}

// NewPubSub initializes the engine
func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[string][]net.Conn),
	}
}

// Subscribe adds a user's raw TCP connection to a channel's broadcast list
func (ps *PubSub) Subscribe(channel string, conn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.channels[channel] = append(ps.channels[channel], conn)
}

// Publish loops through all connections in a channel and sends them the message
func (ps *PubSub) Publish(channel string, message string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subscribers := ps.channels[channel]
	count := 0

	for _, conn := range subscribers {
		// Format a clean message for the listener
		msg := "[MESSAGE - " + channel + "] " + message + "\r\n"
		conn.Write([]byte(msg))
		count++
	}
	
	// Returns the total number of people who received the message
	return count
}

// RemoveClient cleans up dead connections so our server doesn't crash
func (ps *PubSub) RemoveClient(disconnectedConn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for ch, conns := range ps.channels {
		var activeConns []net.Conn
		for _, c := range conns {
			if c != disconnectedConn {
				activeConns = append(activeConns, c) // Keep it if it's still alive
			}
		}
		ps.channels[ch] = activeConns
	}
}