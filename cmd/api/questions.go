package main

import (
	"net/http"

	"github.com/RakibulBh/shaheed-backend/internal/store"
)

type QuestionRequest struct {
	Content  string `json:"content"`
	ParentID *int   `json:"parent_id"`
	Location string `json:"location"`
}

func (app *application) PostQuestion(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtx).(store.User)

	var questionRequest QuestionRequest
	err := app.readJSON(r, &questionRequest)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Convert pointer to int, defaulting to 0 if nil
	var parentID int
	if questionRequest.ParentID != nil {
		parentID = *questionRequest.ParentID
	}

	question, err := app.store.Questions.Create(ctx, user.ID, questionRequest.Content, parentID, questionRequest.Location)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, "success", question)
}

func (app *application) GetQuestions(w http.ResponseWriter, r *http.Request) {

	app.writeJSON(w, http.StatusOK, "success", nil)
}

func (app *application) GetQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateQuestion(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteQuestion(w http.ResponseWriter, r *http.Request) {

}
