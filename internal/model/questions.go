package model

import "time"

type Question struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" validate:"required,min=5"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
