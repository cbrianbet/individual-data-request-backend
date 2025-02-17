package db

import (
	"github.com/palladiumkenya/individual-data-request-backend/internal/config"
	"github.com/palladiumkenya/individual-data-request-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	cfg := config.LoadConfig()
	var DB *gorm.DB
	var err error

	for i := 0; i < 5; i++ { // Retry up to 5 times
		log.Printf("Connecting to database with URL: %s", cfg.DatabaseURL)
		DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to the database")
			return DB, nil
		}

		log.Printf("Failed to connect to database: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Failed to connect to database after multiple attempts: %v", err)
	return nil, err
}

func MigrateDB() (*gorm.DB, error) {
	cfg := config.LoadConfig()

	log.Printf("Migrating database with URL: %s", cfg.DatabaseURL)

	DB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Requesters{}, &models.Requests{}, &models.Assignees{}, &models.Approvals{},
		&models.Approvers{}, &models.Files{}, &models.PointPersons{}, &models.Feedback{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return DB, nil
}
