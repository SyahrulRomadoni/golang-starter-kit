package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"	// Framework web Gin
	"golang-starter-kit/models" // Model database
	"golang-starter-kit/utils"  // Helper untuk (response)
)

func GetRoles(c *gin.Context) {
	var roles []models.Role

	// Mengambil semua role yang belum dihapus (deleted_at IS NULL)
	result := models.DB.Where("deleted_at IS NULL").Find(&roles)

	// Jika terjadi error saat mengambil data, kirim response error
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengambil data", nil))
		return
	}

	// Data berhasil di ambil
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Data berhasil diambil", roles))
}

type CreateRoleInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateRole(c *gin.Context) {
	var input CreateRoleInput
	var role models.Role

	// Input Validation
	ok, resp := utils.InputValidation(c, &input)
	if !ok {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Buat object role baru dengan data dari input
	role = models.Role{
		Name: input.Name,
		CreatedAt: time.Now(),
	}

	// Kondisi Create
	if err := models.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal membuat role", nil))
		return
	}
	
	// Response success
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil dibuat", role))
}

// GetRoleByID menampilkan detail role berdasarkan ID
func GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	
	var role models.Role
	
	// Kondisi data ada atau tidak
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	// Data berhasil ditemukan
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Detail Role", role))
}

type UpdateRoleInput struct {
	Name *string `json:"name"`
}

// UpdateRole mengubah data role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	var input UpdateRoleInput

	// Input Validation
	ok, resp := utils.InputValidation(c, &input)
	if !ok {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Chek Role ada atau tidak
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	// Buat object role baru dengan data dari input
	if input.Name != nil && *input.Name != "" {
		role.Name = *input.Name
	}

	// Kondisi Save
	if err := models.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengupdate role", nil))
		return
	}
	
	// Data berhasil di update
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil diupdate", role))
}

// Delete Role (Soft Delete)
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	
	// Check Role
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	// Field model
	now := time.Now()

	// Kondisi check update deleted_at
	if err := models.DB.Model(&role).Update("deleted_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal menghapus role", nil))
		return
	}

	// Berhasil di delete
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil dihapus", nil))
}
