package api

import (
	"encoding/json"
	"net/http"

	"github.com/garkashy/smart-search/internal/authentication"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&registerRequest)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if the username already exists
	_, err = authentication.GetUserID(r.Context(), s.DB, registerRequest.Username)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	err = authentication.InsertUser(r.Context(), s.DB, registerRequest.Username, registerRequest.Password)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
