package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

func main() {
	//TODO: config
	var connectionString, migrationsPath, migrationsTable string

	flag.StringVar(&connectionString, "connection-string", "", "connection string")
	flag.StringVar(&migrationsPath, "migrations-path", "", "migrations path")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "migrations table")
	flag.Parse()

	if connectionString == "" {
		log.Fatal("connection string is required")
	}
	if migrationsPath == "" {
		log.Fatal("migrations path is required")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}(db)

	goose.SetTableName(migrationsTable)

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations successfully applied")
}
