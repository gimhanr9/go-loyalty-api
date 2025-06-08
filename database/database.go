// database/database.go
package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"go-loyalty-api/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.User{}, &models.LoyaltyTransaction{})
	if err != nil {
		log.Fatal("Failed to migrate models")
	}
}