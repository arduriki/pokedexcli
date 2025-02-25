package pokecache

import (
	"sync"
	"time"
)

// cacheEntry represents a single entry in the cache
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache stores the cached data with thread-safe access
type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex // protects the map from concurrent access
}

// NewCache creates a new cache with a cleanup routine that runs at the given interval
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}

	// Start a goroutine to clean up old entries periodically
	go c.readLoop(interval)

	return c
}

// Add adds a new item to the cache with the given key
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves an item from the cache by key
// Returns the value and a boolean indicating if the key was found
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

// reapLoop runs on a timer and removes old entries from the cache
func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.reap(time.Now().Add(-interval))
	}
}

// reap removes entries older than the given time
func (c *Cache) reap(oldestTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		if entry.createdAt.Before(oldestTime) {
			delete(c.entries, key)
		}
	}
}
