package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	pgvector "github.com/pgvector/pgvector-go"
)

type Chunk struct {
	DocumentID int
	ChunkIndex int
	Content    string
	Embedding  []float32
}

func InsertChunk(ctx context.Context, db *pgxpool.Pool, chunk Chunk) error {
	_, err := db.Exec(ctx, "INSERT INTO chunks (document_id, chunk_index, content) VALUES ($1, $2, $3)",
		chunk.DocumentID, chunk.ChunkIndex, chunk.Content)
	return err
}

func UpdateChunk(ctx context.Context, db *pgxpool.Pool, embedding []float32, chunkIndex int, documentID int) error {
	_, err := db.Exec(ctx, "UPDATE chunks SET embedding = $1 WHERE chunk_index = $2 AND document_id = $3",
		pgvector.NewVector(embedding), chunkIndex, documentID)
	return err
}
