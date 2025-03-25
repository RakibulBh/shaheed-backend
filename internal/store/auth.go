package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	db *sql.DB
}

type RegisterRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (s *AuthStore) Register(ctx context.Context, request RegisterRequest) error {

	query := `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(ctx, query, request.FirstName, request.LastName, request.Email, request.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthStore) HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
