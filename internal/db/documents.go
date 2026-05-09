package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Document struct {
	UserID       string
	FileName     string
	OriginalName string
	FilePath     string
	Size         int64
	Status       string
}

type DocumentListItem struct {
	ID           int       `json:"id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func InsertDocument(ctx context.Context, db *pgxpool.Pool, doc Document) (int, error) {
	var id int
	err := db.QueryRow(ctx,
		"INSERT INTO documents (user_id, file_name, original_name, stored_path, size, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		doc.UserID, doc.FileName, doc.OriginalName, doc.FilePath, doc.Size, doc.Status,
	).Scan(&id)
	return id, err
}

func ListDocuments(ctx context.Context, db *pgxpool.Pool, userID string) ([]DocumentListItem, error) {
	rows, err := db.Query(ctx,
		"SELECT id, file_name, original_name, size, status, created_at FROM documents WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []DocumentListItem
	for rows.Next() {
		var d DocumentListItem
		if err := rows.Scan(&d.ID, &d.FileName, &d.OriginalName, &d.Size, &d.Status, &d.CreatedAt); err != nil {
			continue
		}
		docs = append(docs, d)
	}
	if docs == nil {
		docs = []DocumentListItem{}
	}
	return docs, nil
}
