package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/garkashy/smart-search/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultChunkSize = 200

func FetchEmbeddings(texts []string) ([][]float32, error) {
	body := map[string][]string{"texts": texts}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	mlURL := os.Getenv("ML_SERVICE_URL")
	if mlURL == "" {
		mlURL = "http://localhost:8000"
	}
	resp, err := http.Post(mlURL+"/embed", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get embeddings")
	}

	var response struct {
		Embeddings [][]float32 `json:"embeddings"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Embeddings, nil
}

func splitTextIntoChunks(text string, wordsPerChunk int, documentID int) []db.Chunk {
	words := strings.Fields(text)
	var chunks []db.Chunk

	for i := 0; i < len(words); i += wordsPerChunk {
		end := i + wordsPerChunk
		if end > len(words) {
			end = len(words)
		}
		content := strings.Join(words[i:end], " ")
		chunks = append(chunks, db.Chunk{
			DocumentID: documentID,
			ChunkIndex: len(chunks),
			Content:    content,
		})
	}
	return chunks
}

func updateDocumentStatus(ctx context.Context, pool *pgxpool.Pool, documentID int) error {
	_, err := pool.Exec(ctx, "UPDATE documents SET status = 'READY' WHERE id = $1", documentID)
	return err
}

func ProcessDocument(ctx context.Context, pool *pgxpool.Pool, documentID int) error {

	log.Printf("[ingestion] doc %d: starting", documentID)

	var storedPath string
	err := pool.QueryRow(ctx, "SELECT stored_path FROM documents WHERE id = $1", documentID).Scan(&storedPath)
	if err != nil {
		return fmt.Errorf("fetch document row: %w", err)
	}

	content, err := os.ReadFile(storedPath)
	if err != nil {
		return fmt.Errorf("read file %s: %w", storedPath, err)
	}
	log.Printf("[ingestion] doc %d: read %d bytes from %s", documentID, len(content), filepath.Base(storedPath))

	text := string(content)
	ext := filepath.Ext(storedPath)

	if ext == ".json" {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	chunks := splitTextIntoChunks(text, defaultChunkSize, documentID)
	log.Printf("[ingestion] doc %d: split into %d chunks", documentID, len(chunks))

	for _, chunk := range chunks {
		err = db.InsertChunk(ctx, pool, chunk)
		if err != nil {
			return fmt.Errorf("insert chunk %d: %w", chunk.ChunkIndex, err)
		}
	}
	log.Printf("[ingestion] doc %d: chunks stored in DB", documentID)

	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	embeddings, err := FetchEmbeddings(texts)
	if err != nil {
		return fmt.Errorf("fetch embeddings: %w", err)
	}
	log.Printf("[ingestion] doc %d: received %d embeddings from ML service", documentID, len(embeddings))

	for i, chunk := range chunks {
		err = db.UpdateChunk(ctx, pool, embeddings[i], chunk.ChunkIndex, documentID)
		if err != nil {
			return fmt.Errorf("update chunk %d embedding: %w", chunk.ChunkIndex, err)
		}
	}

	err = updateDocumentStatus(ctx, pool, documentID)
	if err != nil {
		return fmt.Errorf("update status to READY: %w", err)
	}

	log.Printf("[ingestion] doc %d: complete", documentID)
	return nil
}
