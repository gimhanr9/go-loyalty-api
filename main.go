package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-loyalty-api/config"
	"go-loyalty-api/database"
	"go-loyalty-api/routes"
)

func init() {
	// Load environment based on APP_ENV
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
	// Initialize DB
	database.Connect()

	// Init Gin router
	router := gin.Default()

	// Register API routes
	routes.RegisterRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
