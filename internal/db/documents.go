package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Document struct {
	UserID   string
	FileName string
	FilePath string
	Size     int64
	Status   string
}

func InsertDocument(ctx context.Context, db *pgxpool.Pool, doc Document) (int, error) {
	var id int
	err := db.QueryRow(ctx,
		"INSERT INTO documents (user_id, file_name, stored_path, size, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		doc.UserID, doc.FileName, doc.FilePath, doc.Size, doc.Status,
	).Scan(&id)
	return id, err
}
