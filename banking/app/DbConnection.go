package app

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func GetDbClient() *sqlx.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL database: %v\n", err)
	}

	// Configure standard database/sql connection pool settings (previously handled by pgxpool)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL banking database!")
	return db
}