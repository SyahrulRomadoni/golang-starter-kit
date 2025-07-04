// Koneksi ke DB1
package models

import "time"

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	IDRole    uint       `json:"id_role"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`

	// Relation
	Role Role `gorm:"foreignKey:IDRole;references:ID" json:"role"`
}
