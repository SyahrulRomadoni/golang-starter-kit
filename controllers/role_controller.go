package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"  // Framework web Gin
	"golang-starter-kit/config" // Konfigurasi database
	"golang-starter-kit/models" // Model database
	"golang-starter-kit/utils"  // Helper untuk (response)
)

func GetRoles(c *gin.Context) {
	var roles []models.Role

	// Mengambil semua role yang belum dihapus (deleted_at IS NULL)
	result := config.DB.Where("deleted_at IS NULL").Find(&roles)

	// Jika terjadi error saat mengambil data, kirim response error
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.APIResponseError("Gagal mengambil data", nil))
		return
	}

	// Jika tidak ada role yang ditemukan, kirim response kosong
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Data berhasil diambil", roles))
}
