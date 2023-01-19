package models

import "github.com/gofrs/uuid"

type Category struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	CategoryName string    `json:"category_name"`
}
