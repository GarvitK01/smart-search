package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/garkashy/smart-search/internal/authentication"
)

func generateSessionToken() string {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	isValid := authentication.CheckCredentials(r.Context(), s.DB, loginRequest.Username, loginRequest.Password)
	if !isValid {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a new session token
	sessionToken := generateSessionToken()
	if sessionToken == "" {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	// Insert the session token into the database
	err = authentication.InsertSessionToken(r.Context(), s.DB, sessionToken, loginRequest.Username)
	if err != nil {
		http.Error(w, "Failed to insert session token", http.StatusInternalServerError)
		return
	}

	// Set the session token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   86400,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
