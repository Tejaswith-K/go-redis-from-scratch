package storage

import (
	"sync"
	"time"
)

// item wraps our data with an expiration timestamp
type item struct {
	value      string
	expiration int64 // Unix timestamp. 0 means it never expires.
}

type RedisStore struct {
	mu   sync.RWMutex
	data map[string]item // Our map now stores 'item' structs
}

// NewRedisStore initializes the database and starts the garbage collector
func NewRedisStore() *RedisStore {
	store := &RedisStore{
		data: make(map[string]item),
	}
	// Fire off a background thread for Active Expiry
	go store.startActiveExpiry()
	return store
}

func (s *RedisStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item{value: value, expiration: 0}
}

// Get implements Lazy Expiry
func (s *RedisStore) Get(key string) (string, bool) {
	s.mu.RLock()
	itm, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false
	}

	// If it has an expiration and that time has passed...
	if itm.expiration > 0 && time.Now().Unix() > itm.expiration {
		s.mu.Lock()
		delete(s.data, key) // Delete it immediately
		s.mu.Unlock()
		return "", false
	}

	return itm.value, true
}

// Expire adds a TTL to an existing key
func (s *RedisStore) Expire(key string, ttlSeconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	itm, exists := s.data[key]
	if !exists {
		return false
	}

	itm.expiration = time.Now().Unix() + int64(ttlSeconds)
	s.data[key] = itm
	return true
}

// TTL checks how much time a key has left
func (s *RedisStore) TTL(key string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	itm, exists := s.data[key]
	if !exists {
		return -2 // Redis standard: Key does not exist
	}
	if itm.expiration == 0 {
		return -1 // Redis standard: Key exists but has no expiration
	}

	remaining := itm.expiration - time.Now().Unix()
	if remaining < 0 {
		return -2 // Expired
	}
	return remaining
}

// startActiveExpiry is our background garbage collector
func (s *RedisStore) startActiveExpiry() {
	// Wake up every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now().Unix()
		for k, v := range s.data {
			if v.expiration > 0 && now > v.expiration {
				delete(s.data, k) // Throw out the trash
			}
		}
		s.mu.Unlock()
	}
}