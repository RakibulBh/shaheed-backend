package main

import (
	"net/http"
)

func (app *application) PostQuestion(w http.ResponseWriter, r *http.Request) {

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
