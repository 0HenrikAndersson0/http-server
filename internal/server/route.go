package server

import (
	"fmt"
	"lab/api/internal/config"
	"os"
	"strings"
)

func GetFileServerPath(cfg config.Config) string {
	return fmt.Sprintf("/%s/", cfg.FileServerRootPath)
}

func GetPageFromPath(path string, cfg config.Config) (string, bool) {
	if strings.HasPrefix(path, GetFileServerPath(cfg)) {
		return path[1:], false
	}
	switch path {
	case "/":
		return "Pages/main.html", false
	case "/files":
		return "web/templates/FileList.html", false
	case "/logIn":
		return "web/templates/logIn.html", false
	default:
		return FindPage(path, true)
	}
}

func FindPage(path string, useCache bool) (string, bool) {
	cleanPath := strings.TrimPrefix(path, "/")

	if useCache {
		content, exists := GetFromCache(cleanPath)
		if exists {
			return content, false
		}
	}

	// 1. Check Pages/ (Static content)
	fullPath := "Pages/" + cleanPath + ".html"
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, false
	}

	// 2. Check Pages/Auth/ (Secret content)
	secretPath := "Pages/Auth/" + cleanPath + ".html"
	if _, err := os.Stat(secretPath); err == nil {
		return secretPath, true
	}

	// 3. Check web/templates/ (Core pages)
	templatePath := "web/templates/" + cleanPath + ".html"
	if _, err := os.Stat(templatePath); err == nil {
		return templatePath, false
	}

	return "web/templates/404.html", false
}
