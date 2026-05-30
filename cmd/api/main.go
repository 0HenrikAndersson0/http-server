package main

import (
	"fmt"
	"lab/api/internal/auth"
	"lab/api/internal/config"
	"lab/api/internal/server"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/logIn.html")
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		isValid := auth.FormAuthValidation(r, cfg)
		if !isValid {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		auth.SetJWTCookie(w, cfg)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path, _ := getPathAndQuery(r)

		// If we are in File Server mode, serve the generated list for the root
		if cfg.IsFileServer && path == "/" {
			isAuth := auth.AuthenticateRequest(r, w, cfg)
			if !isAuth {
				return
			}
			fmt.Printf("Serving dynamic file list\n")
			list, err := server.ServeFileList(cfg.FileServerRootPath)
			if err != nil {
				http.Error(w, "Failed to generate file list", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, list)
			return
		}

		// Otherwise, serve standard pages or files
		page, requiresAuth := server.GetPageFromPath(path, cfg)
		if requiresAuth {
			isAuth := auth.AuthenticateRequest(r, w, cfg)
			if !isAuth {
				return
			}
		}
		fmt.Printf("Serving: %s\n", page)
		http.ServeFile(w, r, page)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   "auth_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		fmt.Printf("User logged out\n")
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	})

	fmt.Printf("Starting server on :%d...\n", cfg.Port)
	startServer(cfg, mux)
}

func startServer(cfg config.Config, muxInstance *http.ServeMux) bool {
	addr := ":" + strconv.Itoa(cfg.Port)
	if cfg.CertFile == "" || cfg.CertKey == "" {
		err := http.ListenAndServe(addr, muxInstance)
		if err != nil {
			log.Fatalf("Failed to start server %v", err)
		}
		return true
	}
	err := http.ListenAndServeTLS(addr, cfg.CertFile, cfg.CertKey, muxInstance)
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
	return true
}

func getPathAndQuery(r *http.Request) (string, url.Values) {
	path := r.URL.Path
	query := r.URL.Query()
	return path, query
}
