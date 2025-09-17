package main

import (
	"GoSocial/internal/db"
	"GoSocial/internal/env"
	"GoSocial/internal/store"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15mins"),
		},
		jwt: jwtConfig{
			secret: []byte(env.GetString("JWT_SECRET", "change-this-super-secret-key-in-production")),
			issuer: env.GetString("JWT_ISSUER", "gosocial-api"),
			expiry: time.Duration(env.GetInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("Database connection pool established.")

	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
