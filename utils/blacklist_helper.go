package utils

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// TokenEntry menyimpan informasi token beserta waktu kadaluwarsa
type TokenEntry struct {
	ExpiresAt time.Time `json:"ExpiresAt"`
}

var (
	blacklist = make(map[string]TokenEntry)
	mutex     sync.RWMutex
	filePath  = "blacklist.json"
)

// AddToBlacklist menambahkan token ke dalam blacklist beserta waktu kadaluwarsanya
func AddToBlacklist(token string, expiresAt time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	blacklist[token] = TokenEntry{ExpiresAt: expiresAt}
	saveBlacklistToFile()
}

// IsBlacklisted mengecek apakah token sudah ada di dalam blacklist
func IsBlacklisted(token string) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	entry, exists := blacklist[token]
	if !exists {
		return false
	}

	if time.Now().After(entry.ExpiresAt) {
		mutex.RUnlock()
		mutex.Lock()
		delete(blacklist, token)
		mutex.Unlock()
		mutex.RLock()
		saveBlacklistToFile()
		return false
	}

	return true
}

// GetBlacklistedTokens mengembalikan daftar token yang ada di dalam blacklist
func GetBlacklistedTokens() map[string]TokenEntry {
	mutex.RLock()
	defer mutex.RUnlock()
	return blacklist
}

// ClearBlacklist menghapus semua token dari blacklist dan file
func ClearBlacklist() {
	mutex.Lock()
	defer mutex.Unlock()

	blacklist = make(map[string]TokenEntry) // kosongkan map

	// Hapus file atau kosongkan isinya
	_ = os.WriteFile(filePath, []byte("{}"), 0644)
}

// saveBlacklistToFile menyimpan blacklist ke file JSON
func saveBlacklistToFile() {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	_ = encoder.Encode(blacklist)
}

// loadBlacklistFromFile memuat blacklist dari file JSON
func loadBlacklistFromFile() {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&blacklist)
}

// InitBlacklist memuat blacklist dari file saat aplikasi dimulai
func InitBlacklist() {
	loadBlacklistFromFile()
}
