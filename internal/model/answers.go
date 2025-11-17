package model

import (
	"github.com/google/uuid"
	"time"
)

type Answer struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	QuestionId int       `json:"question_id"`
	UserId     uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Text       string    `json:"text" validate:"required,min=1"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
