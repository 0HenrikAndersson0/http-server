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

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming request: Method:%s Path:%s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Loaded RemoteAddr: %s", r.RemoteAddr)
	})
}

func main() {
	mux := http.NewServeMux()
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	mux.Handle("GET /login", LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/logIn.html")
	})))

	mux.Handle("POST /login", LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValid := auth.FormAuthValidation(r, cfg)
		if !isValid {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		auth.SetJWTCookie(w, cfg)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})))

	mux.Handle("/", LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	mux.Handle("/logout", LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   "auth_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		fmt.Printf("User logged out\n")
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	})))

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
