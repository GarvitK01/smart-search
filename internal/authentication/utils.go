package authentication

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func GetUserID(ctx context.Context, db *pgxpool.Pool, username string) (int, error) {

	var id int
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE username = $1`, username).Scan(&id)

	return id, err
}

func getUserHash(ctx context.Context, db *pgxpool.Pool, username string) (string, error) {
	row := db.QueryRow(ctx, `SELECT password FROM users WHERE username = $1`, username)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CheckCredentials(ctx context.Context, db *pgxpool.Pool, username, password string) bool {
	hash, err := getUserHash(ctx, db, username)
	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func InsertSessionToken(ctx context.Context, db *pgxpool.Pool, sessionToken string, username string) error {

	userID, err := GetUserID(ctx, db, username)
	if err != nil {
		return err
	}

	query := `INSERT INTO sessions (user_id, session_token) VALUES ($1, $2)`
	_, err = db.Exec(ctx, query, userID, sessionToken)
	if err != nil {
		log.Println("Error inserting session token:", err)
		return err
	}

	return nil

}

func SessionLookup(ctx context.Context, db *pgxpool.Pool, sessionToken string) (int, error) {
	var userID int
	err := db.QueryRow(ctx, `SELECT user_id FROM sessions WHERE session_token = $1`, sessionToken).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func InsertUser(ctx context.Context, db *pgxpool.Pool, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `INSERT INTO users (username, password) VALUES ($1, $2)`, username, hash)
	return err
}

func DeleteSessionToken(ctx context.Context, db *pgxpool.Pool, sessionToken string) error {
	_, err := db.Exec(ctx, `DELETE FROM sessions WHERE session_token = $1`, sessionToken)
	if err != nil {
		log.Println("Error deleting session token:", err)
		return err
	}

	return nil
}
