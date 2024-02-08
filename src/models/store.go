package models

type Store struct {
	Model

	UserID      string  `gorm:"index" json:"user_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"image_url"`
}