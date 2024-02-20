package models

/*
Description:

	Represents the model for a product in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	StoreID (string): The ID of the store to which the product belongs. Indexed field for efficient querying.
	Name (string): The name of the product.
	Description (*string): The description of the product. Nullable.
	Price (*float32): The price of the product. Nullable.
	Published (bool): Indicates whether the product is published or not.
	ProductImages ([]ProductImage): Slice of product images associated with the product.

Relations:

	Store: Belongs-to relationship to products. Each product belongs to a store.
	Reviews: One-to-many relationship between products and reviews. Each product can have multiple reviews.
	ProductImages: One-to-many relationship between products and product images. Each product can have multiple images.
*/
type Product struct {
	Model

	StoreID       string         `gorm:"index" json:"store_id"`
	Name          string         `json:"name"`
	Description   *string        `json:"description"`
	Price         *float32       `json:"price"`
	Published     bool           `json:"is_published"`
	Reviews       []Review       `json:"reviews"`
	ProductImages []ProductImage `json:"product_images"`
}

/*
Description:

	Represents the model for a product image in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	ProductID (string): The ID of the product to which the product image belongs. Indexed field for efficient querying.
	Url (string): The URL of the product image.

Relations:

	Product: Belongs-to relation to product. Each product image belongs to a product.
*/
type ProductImage struct {
	Model

	ProductID string `gorm:"index" json:"product_id"`
	Url       string `json:"url"`
}
