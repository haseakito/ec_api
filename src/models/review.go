package models

/*
Description:

	Represents the model for a review in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	ProductID (string): The ID of the product to which the review belongs. Indexed field for efficient querying.
	UserID (string): The ID of the user associated with the store. Indexed field for efficient querying.
	Content (string): The content of this review.

Relations:

	Product: Belongs-to relationship to products. Each review belongs to a product.
*/
type Review struct {
	Model

	ProductID string `gorm:"index" json:"product_id"`
	UserID    string `gorm:"index" json:"user_id"`
	Content   string `json:"content"`
}
