package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

// ConnectDB initializes a connection pool to the PostgreSQL database
func ConnectDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Build the database connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Configure the connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database configuration: %v", err)
	}

	// Disable prepared statement caching
	fmt.Println("RuntimeParams:", config.ConnConfig.RuntimeParams)
	config.ConnConfig.RuntimeParams["prefer_simple_protocol"] = "true"

	// Connect to the database using a connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to establish connection pool: %v", err)
	}
	fmt.Println("Database connection pool established successfully")
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection pool closed")
	}
}
