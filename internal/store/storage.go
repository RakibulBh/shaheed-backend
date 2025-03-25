package store

import (
	"context"
	"database/sql"
	"errors"
)

// Errors
var (
	ErrNotFound = errors.New("not found")
	ErrNoRows   = errors.New("user not found")
	ErrConflict = errors.New("conflict")
	ErrInternal = errors.New("internal server error")
	ErrInvalid  = errors.New("invalid input")
)

type Storage struct {
	Questions interface {
		Create(ctx context.Context, question *Question) error
		// Get(ctx context.Context, id string) (*Question, error)
		// Update(ctx context.Context, question *Question) error
		// Delete(ctx context.Context, id string) error
	}
	Auth interface {
		HashPassword(password string) (string, error)
		Register(ctx context.Context, request RegisterRequest) error
		VerifyPassword(password string, hash string) (bool, error)
	}
	User interface {
		GetUserByEmail(ctx context.Context, email string) (User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Questions: &QuestionStore{db: db},
		Auth:      &AuthStore{db: db},
		User:      &UserStore{db: db},
	}
}
