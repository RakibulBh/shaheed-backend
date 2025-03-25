package store

import (
	"context"
	"database/sql"
	"errors"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	PasswordHash string
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (User, error) {

	query := `
	SELECT password_hash
	FROM users
	WHERE email = $1
	`

	var fecthedUser User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&fecthedUser.PasswordHash)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrNoRows
		default:
			return User{}, err
		}
	}

	return fecthedUser, nil
}
