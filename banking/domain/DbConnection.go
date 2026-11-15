package domain


import (
	"context"
	"log"
	"os"
	"time"


	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDbClient() *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:RudigerDeskillz47@localhost:5432/banking?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create PostgreSQL connection pool: %v\n", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping PostgreSQL database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL banking database!")
	return pool
}
