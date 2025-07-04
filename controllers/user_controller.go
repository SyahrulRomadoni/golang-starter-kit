package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"   // Framework web Gin
	"golang.org/x/crypto/bcrypt" // Untuk hashing password
	"golang-starter-kit/models"  // Model database
	"golang-starter-kit/utils"   // Helper (response, jwt)
)

// GetUsers menampilkan semua user
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := models.DB.Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengambil data user", nil))
		return
	}

	// Bersihkan password
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, utils.APIResponseSuccess("Daftar user", users))
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IDRole   int    `json:"id_role" binding:"required"`
}

// CreateUser membuat user baru
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, utils.APIResponseError("Input tidak valid", err.Error()))
	// 	return
	// }
	if err := c.ShouldBindJSON(&input); err != nil {
		// Cek field mana yang kosong
		if (input.Name == "" || input.Name == "null") &&
			(input.Email == "" || input.Email == "null") &&
			(input.Password == "" || input.Password == "null") &&
			input.IDRole == 0 {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Semua field tidak boleh kosong", nil))
			return
		}
		if input.Name == "" || input.Name == "null" {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Name tidak boleh kosong", nil))
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
		if input.IDRole == 0 {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("ID Role tidak boleh kosong", nil))
			return
		}
		return
	}

	// Validasi karakter password hanya boleh a-z, A-Z, 0-9, @ _ - #
	for _, cPass := range input.Password {
		if !((cPass >= 'a' && cPass <= 'z') ||
			(cPass >= 'A' && cPass <= 'Z') ||
			(cPass >= '0' && cPass <= '9') ||
			cPass == '@' || cPass == '_' || cPass == '-' || cPass == '#') {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Password hanya boleh mengandung huruf, angka, dan simbol ( @_-# )", nil))
			return
		}
	}

	// Cek apakah email sudah terdaftar
	var existing models.User
	if err := models.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.APIResponseError("Email sudah terdaftar", nil))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengenkripsi password", nil))
		return
	}

	user := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		IDRole:    uint(input.IDRole),
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal membuat user", nil))
		return
	}

	// Ambil user beserta role-nya
	models.DB.Preload("Role").First(&user, user.ID)
	user.Password = "" // jangan kirim password ke response
	
	c.JSON(http.StatusOK, utils.APIResponseSuccess("User berhasil dibuat", user))
}

// GetUserByID menampilkan detail user berdasarkan ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := models.DB.Where("deleted_at IS NULL").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("User tidak ditemukan", nil))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Detail user", user))
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	IDRole   *int    `json:"id_role"`
}

// UpdateUser mengubah data user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := models.DB.Where("deleted_at IS NULL").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("User tidak ditemukan", nil))
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIResponseError("Input tidak valid", err.Error()))
		return
	}

	// Hanya update field yang dikirim
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil && *input.Email != "" {
		user.Email = *input.Email
	}
	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengenkripsi password", nil))
			return
		}
		user.Password = string(hashedPassword)
	}
	if input.IDRole != nil {
		user.IDRole = uint(*input.IDRole)
	}

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengupdate user", nil))
		return
	}

	// Ambil user beserta role-nya
	models.DB.Preload("Role").First(&user, user.ID)
	user.Password = "" // jangan kirim password ke response
	
	c.JSON(http.StatusOK, utils.APIResponseSuccess("User berhasil diupdate", user))
}

// Delete User (Soft Delete)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := models.DB.Where("deleted_at IS NULL").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("User tidak ditemukan", nil))
		return
	}

	now := time.Now()
	if err := models.DB.Model(&user).Update("deleted_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal menghapus user", nil))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponseSuccess("User berhasil dihapus (soft delete)", nil))
}