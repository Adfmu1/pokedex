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
	ca.mu.Lock()
	defer ca.mu.Unlock()
	// delete each entry if it exists for > interval 
	for key, val := range ca.entries {
		if time.Since(entry.createdAt) >= interval * time.Seconds {
			delete(ca.entries, key)
		}
	}
}

func main() {



}