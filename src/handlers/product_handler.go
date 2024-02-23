package handlers

import (
	"errors"
	"net/http"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new ProductHandler with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*ProductHandler: A pointer to the newly created ProductHandler instance.
*/
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		db: db,
	}
}

/*
Description:

	Get a specific product with the product id provided. Return nil if no record is found.

HTTP Method:

	GET `/api/v1/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h ProductHandler) GetProduct(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	var product models.Product
	res := h.db.Preload("ProductImages").Take(&product, "id = ?", productId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, product)
}

/*
Description:

	Create a review for a specific product with the product id and based on the data provided in the request payload.

HTTP Method:

	POST `/api/v1/products/:id/reviews`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h ProductHandler) CreateReview(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a store with store id
	var product models.Product
	res := h.db.Take(&product, "id = ?", productId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Parsing request payload and validate the data
	// If there is a problem with the request, throw an error
	var req requests.ReviewCreateRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Validate request data
	// If there is a problem with the request, throw an error
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Instantiate a new review
	review := models.Review{
		ProductID: productId,
		UserID:    req.UserID,
		Content:   req.Content,
	}

	// Create a new review for the product
	if res := h.db.Create(&review); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusCreated, review)
}

/*
Description:

	Get all reviews for a specific product with the product id.

HTTP Method:

	GET `/api/v1/product/:id/reviews`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h ProductHandler) GetReviews(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	// If there is no record, then throw a NotFound error
	var reviews []models.Review
	res := h.db.Where("product_id = ?", productId).Find(&reviews)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, reviews)
}

/*
Description:

	Delete a specific review with the review id.

HTTP Method:

	DELETE `/api/v1/products/:id/reviews/:review_id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h ProductHandler) DeleteReview(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get review id from request
	reviewId := c.Param("review_id")

	// Get a product image with product id and product image id
	// If there is no record, then throw a NotFound error
	var review models.Review
	res := h.db.Find(&review, "id = ? AND product_id = ?", reviewId, productId)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Delete the review record
	if res := h.db.Delete(&review); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the review")
}
