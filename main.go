package main

import (
	"go-crud/internal/config"
	"go-crud/internal/handler"
	"go-crud/internal/migrate"
	"go-crud/internal/repository"
	"go-crud/internal/router"
	"go-crud/pkg/database"
	"log"
	"net/http"

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

	deps := handler.Dependencies{
		UserRepo: repository.NewUserRepository(db),
	}
	handler := handler.NewHandler(deps)
	router := router.NewRouter(handler)

	http.ListenAndServe(":8080", router)
}
