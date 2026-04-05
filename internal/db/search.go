package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	pgvector "github.com/pgvector/pgvector-go"
)

type SearchResult struct {
	DocumentID int     `json:"document_id"`
	ChunkIndex int     `json:"chunk_index"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
}

func SemanticSearch(ctx context.Context, db *pgxpool.Pool, embedding []float32, userID string) ([]SearchResult, error) {

	sql := `
	SELECT c.document_id, c.chunk_index, c.content,
       1 - (c.embedding <=> $1) AS score
		FROM chunks c
		JOIN documents d ON c.document_id = d.id
		WHERE c.embedding IS NOT NULL AND d.user_id = $2
		ORDER BY c.embedding <=> $1
		LIMIT 3
	`

	rows, err := db.Query(ctx, sql, pgvector.NewVector(embedding), userID)
	if err != nil {
		log.Println("Failed to search chunks:", err)
		return nil, err
	}

	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var res SearchResult
		err = rows.Scan(&res.DocumentID, &res.ChunkIndex, &res.Content, &res.Score)
		if err != nil {
			log.Println("Failed to scan chunk:", err)
			continue
		}
		results = append(results, res)
	}

	return results, nil

}

func KeywordSearch(ctx context.Context, db *pgxpool.Pool, query string, userID string) ([]SearchResult, error) {

	sql := `
	SELECT c.document_id, c.chunk_index, c.content,
       ts_rank(to_tsvector('english', c.content), plainto_tsquery('english', $1)) AS score
		FROM chunks c
		JOIN documents d ON c.document_id = d.id
		WHERE to_tsvector('english', c.content) @@ plainto_tsquery('english', $1)
		  AND d.user_id = $2
		ORDER BY score DESC
		LIMIT 3
	`

	rows, err := db.Query(ctx, sql, query, userID)
	if err != nil {
		log.Println("Failed to search chunks:", err)
		return nil, err
	}

	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var res SearchResult
		err = rows.Scan(&res.DocumentID, &res.ChunkIndex, &res.Content, &res.Score)
		if err != nil {
			log.Println("Failed to scan chunk:", err)
			continue
		}
		results = append(results, res)
	}

	return results, nil

}
