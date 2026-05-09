package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, db *pgxpool.Pool) error {

	_, err := db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id),
		file_name VARCHAR(255) NOT NULL,
		stored_path VARCHAR(255) NOT NULL,
		size BIGINT NOT NULL,
		status VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS chunks (
		id SERIAL PRIMARY KEY,
		document_id INT NOT NULL references documents(id),
		chunk_index INT NOT NULL,
		content TEXT NOT NULL,
		embedding vector(384),
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id),
		session_token VARCHAR(255) NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 day',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `ALTER TABLE documents ADD COLUMN IF NOT EXISTS original_name VARCHAR(255) NOT NULL DEFAULT ''`)

	return err
}

func Connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, connString)
}
