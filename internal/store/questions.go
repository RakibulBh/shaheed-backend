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
		INSERT INTO questions (content, location, user_id, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, NULLIF($4, 0), $5, $6)
		RETURNING id, created_at, updated_at
	`

	createdAt := time.Now()

	question := &Question{}

	err := s.db.QueryRowContext(ctx, query, content, location, userID, parentID, createdAt, createdAt).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (s *QuestionStore) GetQuestions(ctx context.Context) ([]Question, error) {

	query := `
		SELECT id, content, location, user_id, created_at, updated_at
		FROM questions
		WHERE parent_id IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	questions := []Question{}

	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Content, &question.Location, &question.UserID, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (s *QuestionStore) Get(ctx context.Context, id int) (*Question, error) {

	query := `
		SELECT id, content, location, user_id, parent_id, created_at, updated_at
		FROM questions
		WHERE id = $1	
	`

	question := &Question{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&question.ID, &question.Content, &question.Location, &question.UserID, &question.ParentID, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return question, nil
}

func (s *QuestionStore) Update(ctx context.Context, question *Question) error {

	updatedAt := time.Now()

	query := `
		UPDATE questions SET content = $1, location = $2, updated_at = $3 WHERE id = $4
	`

	_, err := s.db.ExecContext(ctx, query, question.Content, question.Location, updatedAt, question.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *QuestionStore) Delete(ctx context.Context, id int) error {

	query := `
		DELETE FROM questions WHERE id = $1
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
