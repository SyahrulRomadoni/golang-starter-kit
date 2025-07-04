package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"  // Framework web Gin
	"golang-starter-kit/utils"  // Helper (response, blacklist)
)

// Password yang diizinkan
const secretPassword = "secret123"

// Middleware untuk validasi password dari inputan (query param: password)
func validatePassword(c *gin.Context) bool {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIResponseError("Bad request: gagal membaca input", nil))
		return false
	}
	if req.Password != secretPassword {
		c.JSON(http.StatusUnauthorized, utils.APIResponseError("Unauthorized: Password salah", nil))
		return false
	}
	return true
}

// Check Token black list
func GetBlacklistTokens(c *gin.Context) {
	if !validatePassword(c) {
		return
	}
	blacklist := utils.GetBlacklistedTokens()
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Daftar blacklist token", blacklist))
}

// Clear Token black list
func ClearBlacklistTokens(c *gin.Context) {
	if !validatePassword(c) {
		return
	}
	utils.ClearBlacklist()
	c.JSON(http.StatusOK, utils.APIResponseSuccess("Blacklist token berhasil dikosongkan", nil))
}
