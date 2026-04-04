package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		file_name VARCHAR(255) NOT NULL,
		stored_path VARCHAR(255) NOT NULL,
		size BIGINT NOT NULL,
		status VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS chunks (
		id SERIAL PRIMARY KEY,
		document_id INT NOT NULL references documents(id),
		chunk_index INT NOT NULL,
		content TEXT NOT NULL,
		embedding pgvector(384) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)

	return err
}

func Connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {

	return pgxpool.New(ctx, connString)
}
