package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	FullName  string    `json:"fullname"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Mail      string    `json:"mail" gorm:"uniqueIndex"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
