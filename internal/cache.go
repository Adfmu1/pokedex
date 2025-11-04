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

func (ch *Cache) Add(key string, val []byte) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	// create an entry
	entry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	// add en entry
	ch.entries[key] = entry
}

func (ch *Cache) Get(key string) []byte, bool {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	// check for entry in cache
	val, found := ch.entries[key]
	// error
	if !found {
		return nil, found
	}
	// found
	return val, found
}

func main() {



}