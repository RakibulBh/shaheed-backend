package main

import (
	"errors"
	"net/http"

	"github.com/RakibulBh/shaheed-backend/internal/store"
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
