package ingestion

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func processDocument(ctx context.Context, db *pgxpool.Pool, documentID int) error {
	return nil
}
