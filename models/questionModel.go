package models

import "github.com/gofrs/uuid"

type Question struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Score         uint      `json:"score"`
	IsHidden      bool      `json:"is_hidden"`
	Flag          string    `json:"flag"`
	CategoryRefer uuid.UUID `json:"category_id" gorm:"default:null"`
	Category      Category  `gorm:"foreignKey:CategoryRefer"`
}
