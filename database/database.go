// database/database.go
package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"

    "github.com/gimhanr9/go-loyalty-api/models"
)

var DB *gorm.DB

func Connect() {
    db, err := gorm.Open(sqlite.Open("loyalty.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto-migrate the user model
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    DB = db
    log.Println("Connected to database and migrated models.")
}
