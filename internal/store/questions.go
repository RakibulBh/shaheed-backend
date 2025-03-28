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
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	ParentID  int       `json:"parent_id"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *QuestionStore) Create(ctx context.Context, userID int, content string, parentID int, location string) (*Question, error) {

	query := `
		INSERT INTO questions (content, location, user_id, parent_id)
		VALUES ($1, $2, $3, NULLIF($4, 0))
		RETURNING id, created_at, updated_at
	`

	question := &Question{}

	err := s.db.QueryRowContext(ctx, query, content, location, userID, parentID).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return question, nil
}
