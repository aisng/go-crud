package main

import (
	"go-crud/internal/config"
	"go-crud/internal/migrate"
	"go-crud/pkg/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found")
	}
	dbConfig := config.LoadDatabaseConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := migrate.ApplyMigrations(db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	defer db.Close()
}
