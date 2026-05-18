package main

import (
	"log"
	"decent/internal/application/usecase"
	"decent/internal/config"
	"decent/internal/infrastructure/auth"
	"decent/internal/infrastructure/database"
	"decent/internal/infrastructure/http"
	"decent/internal/infrastructure/http/handler"
)

func main() {
	configManager, err := config.NewConfigManager("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := configManager.GetConfig()

	db, err := database.NewPostgresDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	jwtService := auth.NewJWTService(cfg.JWT.Secret)

	bookRepo := db.BookRepository()
	bookUseCase := usecase.NewBookUseCase(bookRepo)

	bookHandler := handler.NewBookHandler(bookUseCase, jwtService)
	authHandler := handler.NewAuthHandler(jwtService)

	router := http.NewRouter(jwtService)
	router.SetupRoutes(bookHandler, authHandler)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(router.Run(cfg.Server.Port))
}