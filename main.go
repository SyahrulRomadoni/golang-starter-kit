package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"  // Untuk memuat variabel dari file .env
	"golang-starter-kit/config" // Package untuk konfigurasi dan koneksi database
	"golang-starter-kit/models" // Package untuk model database (migrasi, dll)
	"golang-starter-kit/routes" // Package untuk routing menggunakan Gin framework
	"golang-starter-kit/utils"  // Helper Blacklist
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Inisialisasi konfigurasi dan koneksi database
	config.ConnectDB()
	// Inisialisasi model database (migrasi, dll)
	models.InitModel()
	// Memuat blacklist dari file
	utils.InitBlacklist()
	// Setup routing menggunakan Gin framework
	r := routes.SetupRoutes()

	// Ambil URL dan PORT dari environment variable
	appURL := os.Getenv("APP_URL")
	appPort := os.Getenv("APP_PORT")
	addr := fmt.Sprintf("%s:%s", appURL, appPort)

	// Jalankan server pada alamat dan port dari .env
	r.Run(addr)
}
