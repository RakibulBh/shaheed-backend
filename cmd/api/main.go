package main

import (
	"log"

	"github.com/RakibulBh/shaheed-backend/internal/db"
	"github.com/RakibulBh/shaheed-backend/internal/env"
	"github.com/RakibulBh/shaheed-backend/internal/store"
)

func main() {

	// TODO: configure the projct

	cfg := config{
		addr:   ":" + env.GetString("PORT", "8080"),
		apiURL: env.GetString("API_URL", "http://localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/shaheed?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10s"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Store
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
