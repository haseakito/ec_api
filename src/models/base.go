package models

import "time"

type Model struct {
	ID        string    `gorm:"primaryKey;size:255;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
