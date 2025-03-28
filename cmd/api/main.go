package main

import (
	"log"
	"time"

	"github.com/RakibulBh/shaheed-backend/internal/db"
	"github.com/RakibulBh/shaheed-backend/internal/env"
	"github.com/RakibulBh/shaheed-backend/internal/redis"
	"github.com/RakibulBh/shaheed-backend/internal/store"
)

func main() {

	cfg := config{
		env:    env.GetString("ENV", "development"),
		addr:   ":" + env.GetString("PORT", "8080"),
		apiURL: env.GetString("API_URL", "http://localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/shaheed?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10s"),
		},
		auth: auth{
			jwtSecret:  env.GetString("AUTH_SECRET", "VERYSECRET"),
			exp:        env.GetDuration("AUTH_EXP", time.Hour*200),
			refreshExp: env.GetDuration("AUTH_REFRESH_EXP", time.Hour*24*7), // 7 days
		},
		llm: llmConfig{
			model:  "gemini-2.0-flash-lite",
			apiKey: env.GetString("GEMINI_API_KEY", "API_KEY_HERE"),
		},
		redis: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			protocol: env.GetInt("REDIS_PROTOCOL", 2),
		},
	}

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Redis
	redis, err := redis.New(cfg.redis.addr, cfg.redis.password, cfg.redis.db, cfg.redis.protocol)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	// Store
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
