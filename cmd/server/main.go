package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/garkashy/smart-search/internal/api"
	"github.com/garkashy/smart-search/internal/db"
)

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

	// Instanstiate the Server
	srv := api.Server{DB: dbPool}
	http.HandleFunc("/upload", srv.UploadHandler)
	http.HandleFunc("/search", srv.SearchHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
