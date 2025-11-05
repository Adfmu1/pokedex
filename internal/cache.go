package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	mu					*sync.RWMutex
	entries				map[string]cacheEntry
}

type cacheEntry struct {
	createdAt			time.Time
	val					[]byte
}

func (ca *Cache) Add(key string, val []byte) {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	// create an entry
	entry := cacheEntry {
		createdAt: time.Now(),
		val: val,
	}
	// add en entry
	ca.entries[key] = entry
}

func (ca *Cache) Get(key string) ([]byte, bool) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	// check for entry in cache
	entry, found := ca.entries[key]
	// error
	if !found {
		return nil, found
	}
	// found
	return entry.val, found
}

func (ca *Cache) reapLoop(interval time.Duration) {
	// ticker to tick after an interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	// check if time has passed
	for range ticker.C {
		ca.mu.Lock()
		// delete entry if it exists for > interval 
		for key, entry := range ca.entries {
			if time.Since(entry.createdAt) >= interval {
				delete(ca.entries, key)
			}
		}
		ca.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		mu:				&sync.RWMutex{},
		entries:		make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}
