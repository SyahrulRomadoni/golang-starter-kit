package middleware

import (
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"     // Framework web Gin
	"github.com/golang-jwt/jwt/v5" // Library JWT untuk parsing dan verifikasi token
	"golang-starter-kit/utils"     // Helper Blacklist
)

// JWTAuth adalah middleware untuk memverifikasi JWT token yang dikirim oleh client
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Authorization header dari request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Jika header tidak ada, tolak permintaan
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Hapus prefix "Bearer " dari token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Jika token tidak diawali dengan "Bearer ", tolak permintaan
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Cek apakah token sudah di-blacklist (misalnya setelah logout)
		if utils.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been logged out"})
			c.Abort()
			return
		}

		// Parse token dan validasi menggunakan JWT secret
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode penandatanganan yang digunakan adalah HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Jika token tidak valid, tolak permintaan
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Lanjutkan ke handler berikutnya jika token valid
		c.Next()
	}
}
