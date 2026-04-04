package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/garkashy/smart-search/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getEmbeddings(chunks []db.Chunk) ([][]float32, error) {

	// Building a JSON body for the ML service FastAPI

	texts := make([]string, len(chunks))

	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	body := map[string]interface{}{
		"texts": texts,
	}

	// Sending the request to the ML service
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://localhost:8000/embed", "application/json", bytes.NewBuffer(jsonBody))
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

func ProcessDocument(ctx context.Context, pool *pgxpool.Pool, documentID int) error {

	var storedPath string
	var docuID int
	err := pool.QueryRow(ctx, "SELECT id, stored_path FROM documents WHERE id = $1", documentID).Scan(&docuID, &storedPath)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(storedPath)
	if err != nil {
		return err
	}

	text := string(content)
	ext := filepath.Ext(storedPath)

	if ext == ".json" {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	chunks := splitTextIntoChunks(text, 500, docuID)

	for _, chunk := range chunks {
		err = db.InsertChunk(ctx, pool, chunk)
		if err != nil {
			return err
		}
	}

	// Send the chunks to the ML service for embeddings
	embeddings, err := getEmbeddings(chunks)
	if err != nil {
		return err
	}

	// Update the chunks with the embeddings
	for i, chunk := range chunks {
		err = db.UpdateChunk(ctx, pool, embeddings[i], chunk.ChunkIndex, docuID)
		if err != nil {
			return err
		}
	}

	return nil
}
