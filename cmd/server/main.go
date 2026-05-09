package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/garkashy/smart-search/internal/api"
	"github.com/garkashy/smart-search/internal/db"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	dbPool, err := db.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	err = db.Migrate(context.Background(), dbPool)
	if err != nil {
		log.Fatal(err)
	}

	srv := api.Server{DB: dbPool}

	mux := http.NewServeMux()
	mux.Handle("/upload", srv.AuthMiddleware(http.HandlerFunc(srv.UploadHandler)))
	mux.Handle("/search", srv.AuthMiddleware(http.HandlerFunc(srv.SearchHandler)))
	mux.Handle("/documents", srv.AuthMiddleware(http.HandlerFunc(srv.DocumentsHandler)))
	mux.Handle("/me", srv.AuthMiddleware(http.HandlerFunc(srv.MeHandler)))
	mux.HandleFunc("/login", srv.LoginHandler)
	mux.HandleFunc("/register", srv.RegisterHandler)
	mux.HandleFunc("/logout", srv.LogoutHandler)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
