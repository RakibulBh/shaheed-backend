package main

import (
	"errors"
	"net/http"
	"strings"
)

func (app *application) PostQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetQuestions(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		app.unauthorizedResponse(w, r, errors.New("no token provided"))
		return
	}

	// Split by space and check if the first part is bearer
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		app.unauthorizedResponse(w, r, errors.New("invalid token"))
		return
	}

	token := parts[1]

	err := app.VerifyToken(token)
	if err != nil {
		app.unauthorizedResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, nil)
}

func (app *application) GetQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteQuestion(w http.ResponseWriter, r *http.Request) {

}
