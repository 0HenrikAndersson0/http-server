package server

import "fmt"

var PageCache = make(map[string]string)

func GetFromCache(path string) (string, bool) {
	content, exists := PageCache[path]
	if exists {
		fmt.Printf("Cache hit for path: %s\n", path)
	} else {
		fmt.Printf("Cache miss for path: %s\n", path)
	}
	return content, exists
}
