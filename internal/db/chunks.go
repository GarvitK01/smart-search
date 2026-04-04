package db

type Chunk struct {

	DocumentID int
	ChunkIndex int
	Content    string
	Embedding  pgvector(384)
	CreatedAt  time.Time
}

func InsertChunk(ctx context.Context, db *pgxpool.Pool, chunk Chunk) error {
	_, err := db.Exec(ctx, "INSERT INTO chunks (document_id, chunk_index, content, embedding) VALUES ($1, $2, $3, $4)", chunk.DocumentID, chunk.ChunkIndex, chunk.Content, chunk.Embedding)
	return err
}