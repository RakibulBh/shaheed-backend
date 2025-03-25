package store

import (
	"context"
	"database/sql"
	"errors"
)

// Errors
var (
	ErrNotFound = errors.New("not found")
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
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Questions: &QuestionStore{db: db},
	}
}
