package api

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/garkashy/smart-search/internal/db"
	"github.com/garkashy/smart-search/internal/ingestion"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

const defaultUserID = "user_01"

func generateSafeFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x%s", b, ext)
}

func validateFile(header *multipart.FileHeader) error {
	if header.Size > maxUploadSize {
		return errors.New("file size exceeds 200MB")
	}

	// Only allow JSON and TXT
	ext := filepath.Ext(header.Filename)
	if ext != ".json" && ext != ".txt" {
		return errors.New("only JSON and TXT files are allowed")
	}

	return nil
}

const maxUploadSize = 200 << 20 // 200 MB

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	uploadedFile, metadata, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}

	defer uploadedFile.Close()

	// Validate the file
	err = validateFile(metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	safeName := generateSafeFilename(metadata.Filename)
	uploadBase := os.Getenv("UPLOAD_DIR")
	if uploadBase == "" {
		uploadBase = "./uploads"
	}
	userDir := filepath.Join(uploadBase, defaultUserID)
	if err = os.MkdirAll(userDir, 0750); err != nil {
		http.Error(w, "Failed to prepare upload directory", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(userDir, safeName)
	newFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, uploadedFile)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	documentID, err := db.InsertDocument(r.Context(), s.DB, db.Document{
		UserID:   defaultUserID,
		FileName: safeName,
		FilePath: filePath,
		Size:     metadata.Size,
		Status:   "PROCESSING",
	})

	if err != nil {
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}

	// Process the document asynchronously
	go func() {
		if err := ingestion.ProcessDocument(context.Background(), s.DB, documentID); err != nil {
			log.Printf("[ingestion] document %d failed: %v", documentID, err)
		}
	}()

	// Return the success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
