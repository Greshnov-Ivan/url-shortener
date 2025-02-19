package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"url-shortener/internal/lib/envparser"
)

func main() {
	connectionString, err := envparser.GetConnectionStringPg()
	if err != nil {
		log.Fatalf("error receiving the connection string: %v", err)
	}

	migrationDirection, exists := os.LookupEnv("MIGRATIONS_DIRECTION")
	if !exists {
		migrationDirection = "up"
	}

	migrationsPath, exists := os.LookupEnv("MIGRATIONS_PATH")
	if !exists {
		migrationsPath = "migrations/"
	}

	migrationsTable, exists := os.LookupEnv("MIGRATIONS_TABLE")
	if !exists {
		migrationsTable = "migrations"
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}()

	goose.SetTableName(migrationsTable)

	switch migrationDirection {
	case "up":
		log.Println("Applying migrations...")
		if err := goose.Up(db, migrationsPath); err != nil {
			log.Fatalf("failed to apply migrations: %v", err)
		}
		log.Println("Migrations applied successfully.")
	case "down":
		log.Println("Rolling back the last migration...")
		if err := goose.Down(db, migrationsPath); err != nil {
			log.Fatalf("failed to rollback migration: %v", err)
		}
		log.Println("Migration rollback completed.")
	default:
		log.Fatalf("Invalid migration direction: %s. Use 'up' or 'down'.", migrationDirection)
	}

	log.Println("migrations successfully applied")
}
