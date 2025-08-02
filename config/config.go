package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func LoadEnv(envFilePath string) error {
	absPath, err := filepath.Abs(envFilePath)
	if err != nil {
		return errors.New("Error getting absolute path for " + envFilePath + " : " + err.Error())
	}

	err = godotenv.Load(absPath)
	if err != nil {
		return errors.New("Error loading .env file at " + absPath)
	}
	return nil
}

func GetDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, errors.New("Failed to connect to database")
	}
	return db, nil
}
