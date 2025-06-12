package database

import (
	"log"

	"github.com/gimhanr9/go-loyalty-api/models"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite" // Register pure Go SQLite driver
)

var DB *gorm.DB

func Connect() {
	var err error

	// Use gorm sqlite driver configured to use pure Go driver
	DB, err = gorm.Open(sqliteDriver.New(sqliteDriver.Config{
		DriverName: "sqlite",
		DSN:        "loyalty.db",
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
