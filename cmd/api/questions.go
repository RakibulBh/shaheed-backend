package main

import (
	"net/http"
	"strconv"

	"github.com/RakibulBh/shaheed-backend/internal/store"
	"github.com/go-chi/chi/v5"
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
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// Convert pointer to int, defaulting to 0 if nil
	var parentID int
	if questionRequest.ParentID != nil {
		parentID = *questionRequest.ParentID
	}

	// Verify if the content should be flagged
	flagged, reason, err := app.store.Questions.VerifyContent(ctx, questionRequest.Content, app.config.llm.model, app.config.llm.apiKey)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if !flagged {
		question, err := app.store.Questions.Create(ctx, user.ID, questionRequest.Content, parentID, questionRequest.Location)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}
		app.writeJSON(w, http.StatusOK, "success", question)
		return
	}

	// Content was flagged so add to the flagged table
	err = app.store.Questions.FlagQuestion(ctx, user.ID, questionRequest.Content, parentID, questionRequest.Location, reason)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusUnprocessableEntity, "question flagged", reason)
}

func (app *application) GetQuestions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	questions, err := app.store.Questions.GetQuestions(ctx)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, "success", questions)
}

func (app *application) GetQuestion(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	ctx := r.Context()

	// Convert the id to an int
	questionID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	question, err := app.store.Questions.Get(ctx, questionID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, "success", question)
}

func (app *application) UpdateQuestion(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	questionID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var questionRequest QuestionRequest
	err = app.readJSON(r, &questionRequest)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	question := &store.Question{
		ID:       questionID,
		Content:  questionRequest.Content,
		Location: questionRequest.Location,
	}
	err = app.store.Questions.Update(ctx, question)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, "success", question)
}

func (app *application) DeleteQuestion(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	ctx := r.Context()

	questionID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Questions.Delete(ctx, questionID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, "success", "question deleted")
}
