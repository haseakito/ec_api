package models

/*
Description:

	Represents the model for a order in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	StoreID (string): The ID of the store to which the order belongs. Indexed field for efficient querying.
	UserID (string): The ID of the user associated with the order. Indexed field for efficient querying.
	OrderItems ([]OrderItem): Slice of order itens associated with the order.
	Paid (bool): Indicates whether the order is paid or not.

Relations:

	Store: Belongs-to relationship to a store. Each order belongs to a store.
	OrderItems: One-to-many relationship between order and order items. Each order can have multiple order items.
*/
type Order struct {
	Model

	StoreID    string      `gorm:"index" json:"store_id"`
	UserID     string      `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
	Paid       bool        `json:"is_paid"`
}

/*
Description:

	Represents the joint model for products and orders in the database.

Fields:

	Model: Embedded struct containing fields for primary key (ID), creation time (CreatedAt), and update time (UpdatedAt).
	ProductID (string): The ID of the product associated with the order. Indexed field for efficient querying.
	OrderID (string): The ID of the order to which the order item belongs. Indexed field for efficient querying.

Relations:

	Product: One-to-one relationship between an order and a product. Each order item has a product.
	Order: One-to-one relationship between an order and an order items. Each order item belongs to an order.
*/
type OrderItem struct {
	Model

	ProductID string  `gorm:"index" json:"product_id"`
	Product   Product `json:"product"`
	OrderID   string  `gorm:"index" json:"order_id"`
}
