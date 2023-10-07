package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func SetupDB() {
	dbUsername := os.Getenv("POSTGRES_USERNAME")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DBNAME")
	dbPortStr := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=Africa/Lagos",
		dbHost,
		dbUsername,
		dbPass,
		dbName,
		dbPortStr,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	err = Migrate()
	if err != nil {
		log.Fatal("Failed to migrate DB:", err)
	}
}
