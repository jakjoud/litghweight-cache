package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Cache struct {
	sync.RWMutex
	items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = value
}

func main() {
	cache := &Cache{
		items: make(map[string]string),
	}

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			key := r.URL.Query().Get("key")
			if val, ok := cache.Get(key); ok {
				w.Write([]byte(val))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if r.Method == http.MethodPost {
			var cacheItem struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
			if err := json.NewDecoder(r.Body).Decode(&cacheItem); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(cacheItem.Value) > 512 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Value too long"))
				return
			}
			cache.Set(cacheItem.Key, cacheItem.Value)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
