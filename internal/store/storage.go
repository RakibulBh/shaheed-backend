package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		Create(ctx context.Context, userID int, content string, parentID int, location string) (*Question, error)
		GetQuestions(ctx context.Context) ([]Question, error)
		Get(ctx context.Context, id int) (*Question, error)
		Update(ctx context.Context, question *Question) error
		Delete(ctx context.Context, id int) error
		VerifyContent(ctx context.Context, content string, modelName string, apiKey string) (bool, string, error)
		FlagQuestion(ctx context.Context, userID int, content string, parentID int, location string, reason string) error
	}
	Auth interface {
		HashPassword(password string) (string, error)
		Register(ctx context.Context, request RegisterRequest) error
		VerifyPassword(password string, hash string) (bool, error)
		GenerateJWT(userID int, expiresAt time.Time, secret string) (string, error)
		VerifyToken(tokenString string, secret string) (*jwt.Token, error)
		StoreRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
		RefreshToken(ctx context.Context, userID int, tokenString string, secret string, refreshExp time.Duration, accessExp time.Duration) (string, string, error)
	}
	User interface {
		GetUserByID(ctx context.Context, id int) (User, error)
		GetUserByEmail(ctx context.Context, email string) (UserData, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Questions: &QuestionStore{db: db},
		Auth:      &AuthStore{db: db},
		User:      &UserStore{db: db},
	}
}
