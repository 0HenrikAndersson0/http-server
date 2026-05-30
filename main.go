package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var template = `<a href="{{filePath}}" class="file-item">
              <div class="file-icon">FILE</div>
              <div class="file-info">
                  <span class="file-name">{{fileName}}</span>
                  <span class="file-meta">{{fileMeta}}</span>
              </div>
              <span class="badge">Asset</span>
          </a>`

var pageCache = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path, _ := getPathAndQuery(r)

		// If we are in File Server mode, serve the generated list for the root
		if config.IsFileServer && path == "/" {
			fmt.Printf("Serving dynamic file list\n")
			fmt.Fprint(w, serveFileList(config.FileServerRootPath))
			return
		}

		// Otherwise, serve standard pages or files
		page := getPageFromPath(path, config)
		fmt.Printf("Serving: %s\n", page)
		http.ServeFile(w, r, page)
	})

	fmt.Printf("Starting server on :%d...\n", config.Port)
	startServer(config, mux)
}

type Config struct {
	Port               int    `json:"port"`
	CertFile           string `json:"certFile"`
	CertKey            string `json:"certKey"`
	IsFileServer       bool   `json:"isFileServer"`
	FileServerRootPath string `json:fileServerRootPath`
}

func startServer(config Config, muxInstance *http.ServeMux) bool {
	addr := ":" + strconv.Itoa(config.Port)
	if config.CertFile == "" || config.CertKey == "" {
		err := http.ListenAndServe(addr, muxInstance)
		if err != nil {
			log.Fatalf("Failed to start server %v", err)
		}
		return true
	}
	err := http.ListenAndServeTLS(addr, config.CertFile, config.CertKey, muxInstance)
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

func getFileServerPath(config Config) string {
	return fmt.Sprintf("/%s/", config.FileServerRootPath)
}

func getPageFromPath(path string, config Config) string {
	if strings.HasPrefix(path, getFileServerPath(config)) {
		return path[1:]
	}
	switch path {
	case "/":
		return "Pages/main.html"
	case "/files":
		return "FileList.html"
	default:
		return findPage(path[1:], true)
	}
}

func serveFileList(fileRootPath string) string {
	data, err := os.ReadFile("FileList.html")
	if err != nil {
		return "<html><body><h1>Error reading FileList.html</h1></body></html>"
	}

	var files []string
	// Ensure the "Files" directory exists or handle it gracefully
	filepath.WalkDir(fileRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	content := string(data)
	// BUG FIX: Assign result back to content
	content = strings.ReplaceAll(content, "{{files}}", createFileListHtml(files))
	return content
}

func createFileListHtml(files []string) string {
	var fileListHtml strings.Builder
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		item := template
		item = strings.ReplaceAll(item, "{{filePath}}", "/"+path)
		item = strings.ReplaceAll(item, "{{fileName}}", info.Name())

		meta := fmt.Sprintf("Size: %d bytes, Updated: %s", info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		item = strings.ReplaceAll(item, "{{fileMeta}}", meta)

		fileListHtml.WriteString(item)
	}
	return fileListHtml.String()
}

func findPage(path string, useCache bool) string {
	if useCache {
		content, exists := getFromCache(path)
		if exists {
			return content
		}
	}

	fullPath := "Pages/" + path + ".html"
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath
	}
	return "Pages/404.html"
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
