package db

import (
	"log"
	"os"

	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err)
	}
	DB, err = gorm.Open(postgres.Open(os.Getenv("dsn")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error loading database", err)
		return
	}
	MigrateErr := DB.AutoMigrate(&models.Account{}, &models.Transations{})
	if MigrateErr != nil {
		log.Fatalf("Migration failed: %v", err)
	}

}
