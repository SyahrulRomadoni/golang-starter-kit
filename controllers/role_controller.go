package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"  // Framework web Gin
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

	// Jika tidak ada role yang ditemukan, kirim response kosong
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Data berhasil diambil", roles))
}

type CreateRoleInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateRolemembuat role baru
func CreateRole(c *gin.Context) {
	var input CreateRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		if input.Name == "" || input.Name == "null" {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Name tidak boleh kosong", nil))
			return
		}
		return
	}

	role := models.Role{
		Name: input.Name,
		CreatedAt: time.Now(),
	}

	if err := models.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal membuat role", nil))
		return
	}
	
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil dibuat", role))
}

// GetRoleByID menampilkan detail role berdasarkan ID
func GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponseSuccess("Detail Role", role))
}

type UpdateRoleInput struct {
	Name *string `json:"name"`
}

// UpdateRole mengubah data role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if input.Name == nil || *input.Name == "" || *input.Name == "null" {
			c.JSON(http.StatusBadRequest, utils.APIResponseError("Name tidak boleh kosong", nil))
			return
		}
		return
	}

	// Hanya update field yang dikirim
	if input.Name != nil && *input.Name != "" {
		role.Name = *input.Name
	}

	if err := models.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengupdate role", nil))
		return
	}
	
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil diupdate", role))
}

// Delete Role (Soft Delete)
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	if err := models.DB.Where("deleted_at IS NULL").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.APIResponseError("Role tidak ditemukan", nil))
		return
	}

	now := time.Now()
	if err := models.DB.Model(&role).Update("deleted_at", &now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal menghapus role", nil))
		return
	}

	c.JSON(http.StatusOK, utils.APIResponseSuccess("Role berhasil dihapus", nil))
}
