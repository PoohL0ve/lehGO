// Package pokecache provides a thread-safe in-memory cache with automatic expiration.
//
// File: cache.go
// Purpose: Stores raw network byte slices indexed by URL string, protecting
// map access with a mutex for safe concurrent access.

package pokecache

import (
	"sync"
	"time"
)

// Holds the raw bytes of an API response and the creation time
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Manages the in-memory map and synchronisation
type Cache struct {
	cache map[string]cacheEntry
	mux   *sync.RWMutex // Protects concurrent reads/writes to the map
}

// Creates a new in-memory cache and starts a background goroutine
// that automatically deletes entries older than the specified interval duration.
func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.RWMutex{},
	}

	// Start the automatic cleanup loop in the background (goroutine)
	go c.reapLoop(interval)

	return c
}

// continuously runs on a ticker interval to remove expired cache entries.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

// checks all entries in the cache and deletes any created before the interval window.
func (c *Cache) reap(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	timeAgo := time.Now().Add(-interval)
	for key, entry := range c.cache {
		if entry.createdAt.Before(timeAgo) {
			delete(c.cache, key)
		}
	}
}

// Inserts a new key-value pair of raw bytes into the cache,
// recording the current timestamp as its creation time.
func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get looks up a key in the cache. It returns the raw byte slice
// and a boolean indicating whether the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}
