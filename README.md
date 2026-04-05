# smart-search

A personal document search engine that combines keyword and semantic search. Upload files, and search them by meaning — not just exact words.

Built from scratch as a learning project. Go backend, PostgreSQL + pgvector for storage and vector search, Python/FastAPI for embeddings. Half buddy-coded with AI agents, half figuring things out the hard way.

## How it works

```
Upload (.txt, .json, .pdf, .docx)
  → Validate & save to disk
  → Store metadata in PostgreSQL
  → Async ingestion:
      Parse → Chunk → Generate embeddings → Store vectors

Search (GET /search?q=your+query)
  → Keyword search (PostgreSQL full-text)    ← runs in parallel
  → Semantic search (pgvector cosine similarity) ←
  → Merge & rank results
  → Return top matches
```

## Architecture

```
cmd/server/main.go          Entry point, DB setup, HTTP routing
internal/api/upload.go       POST /upload — validate, save, trigger ingestion
internal/api/search.go       GET /search  — parallel keyword + semantic search
internal/parser/parser.go    Text extraction (TXT, JSON, PDF, DOCX)
internal/ingestion/worker.go Async pipeline: parse → chunk → embed → store
internal/db/                 PostgreSQL queries, migrations, data models
ml-service/main.py           FastAPI service for sentence-transformer embeddings
```

## Tech stack

- **Go** — API server, ingestion pipeline, file parsing
- **PostgreSQL** — document/chunk metadata, full-text search
- **pgvector** — vector storage and cosine similarity search
- **Python / FastAPI** — ML embedding service (`all-MiniLM-L6-v2`)
- **pdf-xtract** — PDF text extraction (pure Go)
- **go-docx** — DOCX text extraction (pure Go)

## Setup

### Prerequisites

- Go 1.21+
- Python 3.10+
- PostgreSQL 15+ with pgvector extension

### Database

```bash
createdb smartsearch
psql -d smartsearch -c "CREATE EXTENSION IF NOT EXISTS vector"
```

Tables are auto-created on startup via `db.Migrate`.

### ML service

```bash
cd ml-service
python -m venv venv
source venv/bin/activate
pip install fastapi uvicorn sentence-transformers
python main.py
```

The model (`all-MiniLM-L6-v2`) should be pre-downloaded to `ml-service/model/`.

### Run everything

```bash
./start.sh
```

This starts PostgreSQL (if needed), the ML service, and the Go server.

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | — (required) | PostgreSQL connection string |
| `UPLOAD_DIR` | `./uploads` | Where uploaded files are saved |
| `ML_SERVICE_URL` | `http://localhost:8000` | Embedding service URL |

## API

### Upload a file

```bash
curl -X POST http://localhost:8080/upload -F "file=@document.pdf"
```

Accepts `.txt`, `.json`, `.pdf`, `.docx`. Max size: 200MB.

### Search

```bash
curl "http://localhost:8080/search?q=machine+learning"
```

Returns top-ranked chunks combining keyword and semantic similarity scores.

## What I learned

- Go's `net/http` is enough for a clean REST API without frameworks
- pgvector turns PostgreSQL into a capable vector database
- Running keyword and semantic search in parallel with goroutines + `sync.WaitGroup`
- Async ingestion with goroutines and why `context.Background()` matters over `r.Context()`
- Separation of concerns: parser package for format-specific logic, DB package for queries, API package for handlers
