package models

/*
Description:

	Represents the model for a store in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	UserID (string): The ID of the user associated with the store. Indexed field for efficient querying.
	Name (string): The name of the store.
	Description (*string): The description of the store. Nullable.
	ImageUrl (*string): The URL of the store image. Nullable.
	Products ([]Product): Slice of products associated with the store.

Relations:

	Products: One-to-many relationship between stores and products. Each store can have multiple products.
	Order: One-to-many relationship between stores and orders. Each store can have multiple orders.
*/
type Store struct {
	Model

	UserID      string    `gorm:"index" json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	ImageUrl    *string   `json:"image_url"`
	Products    []Product `json:"products"`
	Orders      []Order   `json:"orders"`
}
