package repository

import (
	"log"
	"os"

	"flowdesk/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	log.Println("Connected to PostgreSQL! 🚀")

	
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
	log.Println("Database Migration Completed! ✅")

	return db
}
