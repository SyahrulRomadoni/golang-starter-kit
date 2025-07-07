package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gin-gonic/gin"		// Framework web Gin
	"github.com/golang-jwt/jwt/v5"	// JWT untuk autentikasi
	"golang.org/x/crypto/bcrypt"  	// Untuk hashing password
	"golang-starter-kit/models"   	// Model database
	"golang-starter-kit/utils"    	// Helper (response, jwt, blacklist)
)

// RegisterInput adalah struktur data yang digunakan saat register
type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email,min=6"`
	Password string `json:"password" binding:"required,min=6"`
	IDRole   uint   `json:"id_role" binding:"required"`
}

func Register(c *gin.Context) {
	var user models.User
	var input RegisterInput

	// ------ Validasi ------ //
	// Input Validation
	ok, resp := utils.InputValidation(c, &input)
	if !ok {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Validasi format password (hanya a-z, A-Z, 0-9, @, #, $)
	if !utils.InputValidationPasswordCriteria(input.Password) {
		c.JSON(http.StatusBadRequest, utils.APIResponseError(
			"Password hanya boleh berisi huruf, angka, dan karakter @, #, $", nil))
		return
	}

	// Cek apakah email sudah terdaftar
	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.APIResponseError("Email sudah terdaftar", nil))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengenkripsi password", nil))
		return
	}
	// ------ END Validasi ------ //

	// Buat object user baru dengan data dari input
	user = models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		IDRole:    input.IDRole,
		CreatedAt: time.Now(),
	}

	// Simpan user baru ke database
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal menyimpan data ke database", nil))
		return
	}

	// Kirim response sukses dengan data user yang baru dibuat
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Registrasi berhasil", user))
}

// LoginInput adalah struktur data yang digunakan saat login
type LoginInput struct {
	Email    string `json:"email" binding:"required"`    // Wajib diisi
	Password string `json:"password" binding:"required"` // Wajib diisi
}

func Login(c *gin.Context) {
	// Parsing dan validasi input JSON dari body request
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Cek field mana yang kosong
		if (input.Email == "" || input.Email == "null") && (input.Password == "" || input.Password == "null") {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Email dan Password tidak boleh kosong", nil))
			return
		}
		if input.Email == "" || input.Email == "null" {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Email tidak boleh kosong", nil))
			return
		}
		if input.Password == "" || input.Password == "null" {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Password tidak boleh kosong", nil))
			return
		}
		return
	}

	// Ambil user dari database berdasarkan email, sekaligus preload relasi role-nya
	var user models.User

	// Check Email ada atau tidak
	if err := models.DB.Preload("Role").Where("email = ?", input.Email).Where("deleted_at IS NULL").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, utils.APIResponseError("Email tidak ditemukan", nil))
		return
	}

	// Cek apakah password yang diinput cocok dengan password yang di-hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.APIResponseError("Password yang anda masukan salah", nil))
		return
	}

	// Generate token JWT berdasarkan ID dan email user
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal membuat token", nil))
		return
	}

	// Siapkan data response yang berisi token dan informasi user
	expiredAt := time.Now().Add(time.Hour * 24) // expired JWT 24 jam, sesuaikan jika perlu
	data := gin.H{
		"expired": expiredAt.Format(time.RFC3339),
		"token":   token,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"id_role": user.IDRole,
			"role": gin.H{
				"id":   user.Role.ID,
				"name": user.Role.Name,
			},
		},
	}

	// Kirim response sukses dengan data user dan token
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Login berhasil", data))
}

func Logout(c *gin.Context) {
	// Ambil token dari header Authorization, lalu hapus prefix "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse dan validasi token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Jika token tidak valid, kirim response error
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.APIResponseError("Token tidak valid", nil))
		return
	}

	// Ambil waktu kadaluarsa dari token, dan tambahkan ke blacklist
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			// Konversi timestamp ke waktu dan tambahkan ke daftar blacklist
			expTime := time.Unix(int64(exp), 0)
			utils.AddToBlacklist(tokenString, expTime)
		}
	}

	// Kirim response logout sukses
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Berhasil logout", nil))
}
