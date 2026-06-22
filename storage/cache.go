package storage

import "sync"

type RedisStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewRedisStore initializes an instance of our thread-safe memory vault
func NewRedisStore() *RedisStore {
	return &RedisStore{
		data: make(map[string]string),
	}
}

// Set safely locks the map and updates the key
func (s *RedisStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get safely locks the map for reading and fetches the value
func (s *RedisStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[key]
	return val, exists
}