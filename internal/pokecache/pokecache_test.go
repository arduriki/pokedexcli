package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	key := "test-key"
	val := []byte("test-value")

	// Add item to cache
	cache.Add(key, val)

	// Retrieve item from cache
	retrievedVal, ok := cache.Get(key)
	if !ok {
		t.Errorf("Expected to find key %s in cache, bit it was not found", key)
	}

	// Compare values
	if string(retrievedVal) != string(val) {
		t.Errorf("Expected value %s, got %s", string(val), string(retrievedVal))
	}
}

func TestCacheReap(t *testing.T) {
	// Create cache with very short reap interval for testing
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	key := "test-key"
	val := []byte("test-value")

	// Add item to cache
	cache.Add(key, val)

	// Verify item exists
	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("Expected to find key %s in cache, but it was not found", key)
	}

	// Wait for longer than the reap interval
	time.Sleep(interval * 2)

	// Verify item is now gone
	_, ok = cache.Get(key)
	if ok {
		t.Errorf("Expected key %s to be reaped from cache, but it was found", key)
	}
}
