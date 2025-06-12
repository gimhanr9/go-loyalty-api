package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
