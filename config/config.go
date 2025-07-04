package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB // Koneksi ke DB1
)

// LoadEnv memuat file .env ke dalam variabel lingkungan
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// ConnectDB menghubungkan aplikasi ke dua database: PostgreSQL dan MySQL
func ConnectDB() {
	LoadEnv()
	var err error

	// ======= Database ======= //
	// DB DSN
	database := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	// Koneksi DB
	DB, err = gorm.Open(postgres.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to BD: ", err)
	}
	sqlDBPg, _ := DB.DB()
	sqlDBPg.SetMaxIdleConns(10)
	sqlDBPg.SetConnMaxLifetime(time.Hour)
}
