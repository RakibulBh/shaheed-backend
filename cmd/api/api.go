package main

import (
	"net/http"
	"time"

	"github.com/RakibulBh/shaheed-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Route("/health", func(r chi.Router) {
		r.Get("/", app.Healthcheck)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/questions", func(r chi.Router) {
			r.Post("/", app.PostQuestion)
			r.Get("/", app.GetQuestions)
			r.Get("/{id}", app.GetQuestion)
			r.Put("/{id}", app.UpdateQuestion)
			r.Delete("/{id}", app.DeleteQuestion)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.Register)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return srv.ListenAndServe()
}
