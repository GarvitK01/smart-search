package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Chunk struct {
	DocumentID int
	ChunkIndex int
	Content    string
}

func InsertChunk(ctx context.Context, db *pgxpool.Pool, chunk Chunk) error {
	_, err := db.Exec(ctx, "INSERT INTO chunks (document_id, chunk_index, content) VALUES ($1, $2, $3)",
		chunk.DocumentID, chunk.ChunkIndex, chunk.Content)
	return err
}
