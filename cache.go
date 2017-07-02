package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)


// CacheEvent
type CacheEvent struct {
	Key string
	File string
	Op CacheOp
}

type CacheOp uint8

const (
	UPDATE = 1
	DELETE = 2
)

func (event CacheEvent) String() string {
	switch (event.Op) {
		case UPDATE:
			return fmt.Sprintf("[Event] UPDATE \n\tFile: %s\n\tKey: %s", event.File, event.Key) 
		case DELETE:
			return fmt.Sprintf("[Event] DELETE \n\tFile: %s\n\tKey: %s", event.File, event.Key) 
	}

	return ""
} 
// /CacheEvent

type CacheItem struct {
	Key string 
	File string 
}

/* wrap items in concurrency safe lock */
type itemIndex struct {
	sync.RWMutex
	items map[string]string
}

type Cache struct {
	Events chan CacheEvent
	items *itemIndex
	done chan bool
}

func NewCache() *Cache {
	cache := Cache{}

	cache.done = make(chan bool, 1)

	cache.Listen()

	return &cache
}

func (cache *Cache) Close() {
	log.Println("Closing cache...")
	cache.done <- true
	close(cache.Events)
}

/*
 * Listen for cache events
 */
func (cache *Cache) Listen() {
	cache.Events = make(chan CacheEvent)

	cache.items = &itemIndex{items: make(map[string]string)}

	go func() {
		for {
			select {
			case <- cache.done:

				return

			case event := <- cache.Events:

				log.Println("Event received...")

				if event.Op == UPDATE {
					cache.updateCache(event.File, event.Key)
				}

				if event.Op == DELETE {
					if event.File != "" {
						cache.deleteCacheItemByFile(event.File)
					} else {
						cache.deleteCacheItem(event.Key)
					}
				}
			}
		}
	}()

}

// set the item in the index
func (cache *Cache) updateCache(file string, key string) {
	log.Println("Update: ", file, " ", key)

	cache.items.Lock()
	cache.items.items[key] = file
	cache.items.Unlock()
}

func (cache *Cache) deleteCacheItemByFile(file string) {
	log.Println("deleting by filename: " + file)

	cache.items.Lock()

	defer cache.items.Unlock()

	for key := range cache.items.items {
		if cache.items.items[key] == file {
			delete(cache.items.items, key)
			return
		}
	}
}

// delete the file and remove from index
func (cache *Cache) deleteCacheItem(key string) {
	log.Println("REMOVE: ", key)

	file, exists := cache.items.items[key]

	if !exists {
		log.Println("Key not found in cache: " + key)
		log.Println(cache.items.items)
		return
	}

	cache.items.Lock()
	delete(cache.items.items, key)
	cache.items.Unlock()

	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		log.Println("File does not exist")
	} else {
		os.Remove(file)
	}

}
