package main

import (
	"net/http"

	"github.com/RakibulBh/shaheed-backend/internal/store"
)

func (app *application) PostQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetQuestions(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtx).(store.User)

	app.writeJSON(w, http.StatusOK, "success", user)
}

func (app *application) GetQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteQuestion(w http.ResponseWriter, r *http.Request) {

}
