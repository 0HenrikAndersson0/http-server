package server

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Template = `<a href="{{filePath}}" class="file-item">
              <div class="file-icon">FILE</div>
              <div class="file-info">
                  <span class="file-name">{{fileName}}</span>
                  <span class="file-meta">{{fileMeta}}</span>
              </div>
              <span class="badge">Asset</span>
          </a>`

func ServeFileList(fileRootPath string) (string, error) {
	data, err := os.ReadFile("web/templates/FileList.html")
	if err != nil {
		return "", err
	}

	var files []string
	filepath.WalkDir(fileRootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	content := string(data)
	content = strings.ReplaceAll(content, "{{files}}", createFileListHtml(files))
	return content, nil
}

func createFileListHtml(files []string) string {
	var fileListHtml strings.Builder
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("Failed getting file stat: %v", err)
			continue
		}

		item := Template
		item = strings.ReplaceAll(item, "{{filePath}}", "/"+path)
		item = strings.ReplaceAll(item, "{{fileName}}", info.Name())

		meta := fmt.Sprintf("Size: %d bytes, Updated: %s", info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		item = strings.ReplaceAll(item, "{{fileMeta}}", meta)

		fileListHtml.WriteString(item)
	}
	return fileListHtml.String()
}
