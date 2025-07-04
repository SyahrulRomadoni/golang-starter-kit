package models

import (
	"golang-starter-kit/config" // Mengimpor koneksi database dari package config
	"gorm.io/gorm"              // ORM (Object Relational Mapper) dari GORM
)

// Variabel global DB agar bisa diakses dari file lain
var DB *gorm.DB

// InitModel digunakan untuk menginisialisasi model dan menghubungkan DB dari config
func InitModel() {
	// Ambil koneksi database dari package config dan simpan ke variabel global DB
	DB = config.DB
}
