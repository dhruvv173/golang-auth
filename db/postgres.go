package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB()error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := getConnectionString()
	db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Cant initialize DB:", err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Cannot run migrations:", err)
	}
	return nil
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting underlying database connection:", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Error closing database connection:", err)
	}
}

func getConnectionString() string {
	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB URL NOT FOUND")
	}
	return dbURL
}

func GetDB() *gorm.DB{
	return db
}