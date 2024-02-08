package database

import (
	"log"
	"os"

	"github.com/haseakito/ec_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 
func Init() *gorm.DB {
	// Get the database connection string from .env
	dsn := os.Getenv("DB_URL")

	// Try to establish connection to database 
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// If the connection was unsuccessful, then throw an error
	if err != nil {
		log.Fatal(err)
	}

	// Add uuid-ossp extension to PostgreSQL
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Run migration
	db.AutoMigrate(&models.Store{})

	return db
}
