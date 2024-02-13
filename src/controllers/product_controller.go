package controllers

import (
	"errors"
	"net/http"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/haseakito/ec_api/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new ProductController with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*ProductController: A pointer to the newly created ProductController instance.
*/
func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
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
func (pc ProductController) GetProduct(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	var product models.Product
	res := pc.db.Preload("ProductImages").Take(&product, "id = ?", productId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, product)
}

/*
Description:

	Update a specific product with the product id provided and based on the data in the request payload.

HTTP Method:

	PATCH `/api/v1/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (pc ProductController) UpdateProduct(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	var product models.Product
	res := pc.db.First(&product, "id = ?", productId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Parsing request payload and validate the data
	// If there is a problem with the request, throw an error
	var req requests.ProductUpdateRequest
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

	// Update product fields if the fields are not empty
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = &req.Description
	}
	if req.Price != 0 {
		product.Price = &req.Price
	}
	if req.Published != false {
		product.Published = req.Published
	}

	// Update product with data
	// If the update is unsuccessful, then throw an error
	if res := pc.db.Save(&product); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, product)
}

/*
Description:

	Upload products image files to storage and update the product with the product id.

HTTP Method:

	POST `/api/v1/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (pc ProductController) UploadImages(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	// If there is no record, then throw a NotFound error
	var product models.Product
	res := pc.db.First(&product, "id = ?", productId)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Get files from request
	files := form.File["images"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, "No files uploaded")
		return nil
	}

	// Iterate over images and upload them to S3
	for _, file := range files {
		// Validate request file
		// If there is a problem with the request, throw an error
		if err := requests.ValidateFile(file); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return nil
		}

		// Upload file to AWS S3 bucket
		// If the upload is unsuccessful, then throw an error
		url, err := utils.Upload(file, "products/")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}

		// Instantiate a new product image
		product_image := models.ProductImage{
			Url:       url,
			ProductID: productId,
		}

		// Create a new product image
		if res := pc.db.Create(&product_image); res.Error != nil {
			c.JSON(http.StatusInternalServerError, res.Error)
			return nil
		}
	}

	return c.JSON(http.StatusOK, "Successfully uploaded the images")
}

/*
Description:

	Delete a specific product with the product id and delete the corresponding object in storage.

HTTP Method:

	DELETE `/api/v1/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (pc ProductController) DeleteProduct(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	var product models.Product
	res := pc.db.Preload("ProductImages").Find(&product, "id = ?", productId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// If the store has image urls, delete the corresponding objects from S3
	if len(product.ProductImages) > 0 {
		// Iterate over product images and delete the corresponding object from S3
		for _, image := range product.ProductImages {
			// Delete the corresponsing object from S3
			if err := utils.Delete(image.Url); err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return nil
			}

			// Delete the product image record
			if res := pc.db.Delete(&image); res.Error != nil {
				c.JSON(http.StatusInternalServerError, res.Error)
				return nil
			}
		}
	}

	// Delete a product
	// If the delete is unsuccessful, then throw an error
	if res := pc.db.Delete(&product, "id = ?", productId); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the product")
}

/*
Description:

	Delete a specific product image with the product_image id and delete the object in storage.

HTTP Method:

	DELETE `/api/v1/products/:id/assets/:image_id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (pc ProductController) DeleteProductImage(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get product image id from request
	productImageId := c.Param("image_id")

	// Get a product image with product id and product image id
	// If there is no record, then throw a NotFound error
	var productImage models.ProductImage
	res := pc.db.First(&productImage, "id = ? AND product_id = ?", productImageId, productId)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Delete the corresponsing object from S3
	if err := utils.Delete(productImage.Url); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}

	// Delete the product image record
	if res := pc.db.Delete(&productImage); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the product image")
}
