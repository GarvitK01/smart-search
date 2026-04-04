package ingestion

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/garkashy/smart-search/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	return nil
}
