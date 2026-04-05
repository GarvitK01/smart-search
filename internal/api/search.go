package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/garkashy/smart-search/internal/db"
	"github.com/garkashy/smart-search/internal/ingestion"
)

const maxResults = 3

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	userID := defaultUserID

	var (
		keywordResults  []db.SearchResult
		semanticResults []db.SearchResult
		keywordErr      error
		semanticErr     error
		wg              sync.WaitGroup
	)

	// Goroutine A: keyword search (only needs the raw query string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		keywordResults, keywordErr = db.KeywordSearch(r.Context(), s.DB, query, userID)
	}()

	// Goroutine B: fetch embedding then semantic search
	wg.Add(1)
	go func() {
		defer wg.Done()
		embeddings, err := ingestion.FetchEmbeddings([]string{query})
		if err != nil {
			semanticErr = err
			return
		}
		semanticResults, semanticErr = db.SemanticSearch(r.Context(), s.DB, embeddings[0], userID)
	}()

	wg.Wait()

	if keywordErr != nil {
		http.Error(w, "Failed keyword search", http.StatusInternalServerError)
		return
	}
	if semanticErr != nil {
		http.Error(w, "Failed semantic search", http.StatusInternalServerError)
		return
	}

	merged := mergeResults(keywordResults, semanticResults)

	// 5. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(merged)
}

func mergeResults(keyword, semantic []db.SearchResult) []db.SearchResult {

	// Deduplicate Result by DocumentID + ChunkIndex
	deduped := make(map[string]db.SearchResult)
	for _, result := range keyword {
		key := fmt.Sprintf("%d-%d", result.DocumentID, result.ChunkIndex)
		deduped[key] = result
	}
	for _, result := range semantic {
		key := fmt.Sprintf("%d-%d", result.DocumentID, result.ChunkIndex)
		if existing, ok := deduped[key]; ok {
			existing.Score += result.Score
			deduped[key] = existing
		} else {
			deduped[key] = result
		}
	}

	// Convert back to slice
	var results []db.SearchResult
	for _, result := range deduped {
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}
