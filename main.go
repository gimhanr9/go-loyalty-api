package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gimhanr9/go-loyalty-api/config"
	"github.com/gimhanr9/go-loyalty-api/database"
	"github.com/gimhanr9/go-loyalty-api/routes"
)

func init() {
	env := os.Getenv("APP_ENV")
	var err error

	if env == "production" {
		err = godotenv.Load(".env.production")
		log.Println("Loaded .env.production")
	} else {
		err = godotenv.Load(".env.development")
		log.Println("Loaded .env.development")
	}

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	database.Connect()

	router := gin.Default()

	routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
