package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/RakibulBh/shaheed-backend/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {

	// Parse the request
	var payload RegisterRequest
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// Check if email exists already
	_, err = app.store.User.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			// do nothing
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	// Validate the length of each field
	if len(payload.FirstName) < 2 || len(payload.LastName) < 2 || len(payload.Email) < 5 || len(payload.Password) < 8 || payload.Password != payload.PasswordConfirm {
		app.badRequestResponse(w, r, errors.New("invalid request payload"))
		return
	}

	// validate password matches
	if payload.Password != payload.PasswordConfirm {
		app.badRequestResponse(w, r, errors.New("password does not match"))
		return
	}

	// Hash the passowrd
	hash, err := app.store.Auth.HashPassword(payload.Password)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	// store the user in the database
	err = app.store.Auth.Register(ctx, store.RegisterRequest{FirstName: payload.FirstName, LastName: payload.LastName, Email: payload.Email, PasswordHash: hash})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, nil)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {

	// Parse the request
	var payload LoginRequest
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// Fetch user from the database
	user, err := app.store.User.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNoRows):
			app.badRequestResponse(w, r, errors.New("invalid credentials"))
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	// Verify password matches with the database hash
	passwordMatches, err := app.store.Auth.VerifyPassword(payload.Password, user.PasswordHash)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if !passwordMatches {
		app.badRequestResponse(w, r, errors.New("invalid credentials"))
		return
	}

	// Generate a JWT token
	token, err := app.GenerateJWT(user.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]string{
		"token": token,
	})
}

func (app *application) GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(app.config.auth.exp).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.config.auth.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (app *application) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(app.config.auth.secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
