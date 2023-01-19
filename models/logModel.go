package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Log struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	UserRefer     uuid.UUID `json:"user_id"`
	User          User      `gorm:"foreignKey:UserRefer"`
	QuestionRefer uuid.UUID `json:"question_id"`
	Question      Question  `gorm:"foreignKey:QuestionRefer"`
	Status        bool      `json:"status"`
	Date          time.Time `json:"date"`
}
