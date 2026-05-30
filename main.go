package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var pageCache = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path, _ := getPathAndQuery(r)
		page := getPageFromPath(path)
		fmt.Printf("Serving page: %s\n", page)
		http.ServeFile(w, r, page)
	})
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	fmt.Printf("Starting server on :%d...", config.Port)
	startServer(config, mux)
}

type Config struct {
	Port     int    `json:"port"`
	CertFile string `json:"certFile"`
	CertKey  string `json:"certKey"`
}

func startServer(config Config, muxInstance *http.ServeMux) bool {
	if config.CertFile == "" || config.CertKey == "" {
		err := http.ListenAndServe(":"+strconv.Itoa(config.Port), muxInstance)
		if err != nil {
			log.Fatalf("Failed to start server %v", err)
		}
		return true
	}
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.CertFile, config.CertKey, muxInstance)
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
	return true
}

func readConfig() (Config, error) {
	data, err := os.ReadFile("conf.json")
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func getPathAndQuery(r *http.Request) (string, url.Values) {
	path := r.URL.Path
	query := r.URL.Query()
	return path, query
}

func getQueryAsString(r *http.Request) string {
	return r.URL.Query().Encode()
}

func getPageFromPath(path string) string {
	switch path {
	case "/":
		return "Pages/main.html"
	default:
		return findPage(path[1:], true)
	}
}

func findPage(path string, useCache bool) string {
	if useCache {
		content, exists := getFromCache(path)
		if exists {
			return content
		}
	}

	_, err := os.Stat("Pages/" + path + ".html")
	if err == nil {
		return "Pages/" + path + ".html"
	} else if os.IsNotExist(err) {
		return "Pages/404.html"
	} else {
		fmt.Println("error checking file:", err)
		return "Pages/404.html"
	}
}

func addToCache(path string, content string) {
	pageCache[path] = content
}

func getFromCache(path string) (string, bool) {
	content, exists := pageCache[path]
	if exists {
		fmt.Printf("Cache hit for path: %s\n", path)
	} else {
		fmt.Printf("Cache miss for path: %s\n", path)
	}
	return content, exists
}
