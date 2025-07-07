package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"   				// Framework web Gin
	"golang.org/x/crypto/bcrypt" 				// Untuk hashing password
	"golang-starter-kit/models"  				// Model database
	"golang-starter-kit/utils"   				// Helper (response, jwt)
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
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email,min=6"`
	Password string `json:"password" binding:"required,min=6"`
	IDRole   int    `json:"id_role" binding:"required"`
}

// CreateUser membuat user baru
func CreateUser(c *gin.Context) {
	var user models.User
	var input CreateUserInput

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

	// Buat user baru
	user = models.User{
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

// Struct untuk input update user
type UpdateUserInput struct {
	Name     *string `json:"name" binding:"omitempty,min=3"`
	Email    *string `json:"email" binding:"omitempty,email,min=6"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	IDRole   *int    `json:"id_role"`
}

// UpdateUser mengubah data user
func UpdateUser(c *gin.Context) {
	
	id := c.Param("id")

	var user models.User
	var input UpdateUserInput
	
	// ------ Validasi Input JSON ------ //
	// Validation Input
	ok, resp := utils.InputValidation(c, &input)
	if !ok {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Validasi format password (hanya a-z, A-Z, 0-9, @, #, $)
	if !utils.InputValidationPasswordCriteria(*input.Password) {
		c.JSON(http.StatusBadRequest, utils.APIResponseError(
			"Password hanya boleh berisi huruf, angka, dan karakter @, #, $", nil))
		return
	}

	// Cek apakah email sudah terdaftar
	if input.Email != nil && *input.Email != "" {
		var existing models.User
		if err := models.DB.Where("email = ? AND id != ?", *input.Email, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Email sudah terdaftar", nil))
			return
		}
	}

	// Check User
	if err := models.DB.Where("deleted_at IS NULL").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("User tidak ditemukan", nil))
		return
	}
	// ------ END Validasi Input JSON ------ //

	// Update data yang dikirim
	if input.Name != nil && *input.Name != "" {
		user.Name = *input.Name
	}
	if input.Email != nil && *input.Email != "" {
		user.Email = *input.Email
	}
	if input.Password != nil && *input.Password != "" {
		// Hash password
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

	// Preload Role dan hilangkan password dari output
	models.DB.Preload("Role").First(&user, user.ID)
	user.Password = ""

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

	c.JSON(http.StatusOK, utils.APIResponseSuccess("User berhasil dihapus", nil))
}