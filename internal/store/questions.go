package store

import (
	"context"
	"database/sql"
	"time"
)

type QuestionStore struct {
	db *sql.DB
}

type Question struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *QuestionStore) Create(ctx context.Context, question *Question) error {
	return nil
}
