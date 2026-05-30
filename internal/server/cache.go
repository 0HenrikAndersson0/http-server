package server

import (
	"fmt"
	"sync"
)

var PageCache sync.Map

func GetFromCache(path string) (string, bool) {
	content, exists := PageCache.Load(path)
	if exists {
		fmt.Printf("Cache hit for path: %s\n", path)
	} else {
		fmt.Printf("Cache miss for path: %s\n", path)
	}
	return content.(string), exists
}

func WriteToCache(path string) {
	PageCache.Store(path, path)
}
