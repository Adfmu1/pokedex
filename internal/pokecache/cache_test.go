package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const internal = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "google.com",
			val: []byte("google_data"),
		},
		{
			key: "amazon.com",
			val: []byte("amazon_data"),
		},
		{
			key: "badsite.com",
			val: []byte("badsite_data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(int(internal))
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Error("Expected to find a key")
				return
			}
			if string(val) != string(c.val) {
				t.Error("Expected to find a value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 1
	const waitTime = baseTime + 1*time.Second
	cache := NewCache(baseTime)
	cache.Add("google.com", []byte("google_data"))

	if _, ok := cache.Get("google.com"); !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get("google.com"); ok {
		t.Errorf("expected to not find key")
		return
	}
}
