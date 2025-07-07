package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	// Panjang secret: 32 byte
	secretLength := 32
	secret := make([]byte, secretLength)

	// Generate random byte
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Println("Gagal generate secret:", err)
		return
	}

	// Encode ke Base64
	encoded := base64.StdEncoding.EncodeToString(secret)

	// Tampilkan ke terminal
	fmt.Printf("JWT_SECRET=%s\n", encoded)
}
