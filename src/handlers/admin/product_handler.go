package admin

import (
	"net/http"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/haseakito/ec_api/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminProductHandler struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new AdminProductHandler with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*AdminProductHandler: A pointer to the newly created AdminProductHandler instance.
*/
func NewAdminProductHandler(db *gorm.DB) *AdminProductHandler {
	return &AdminProductHandler{
		db: db,
	}
}

/*
Description:

	Update a specific product with the product id provided and based on the data in the request payload.

HTTP Method:

	PATCH `/api/v1/admin/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminProductHandler) UpdateProduct(c echo.Context) error {
	// Get product id from request
	productID := c.Param("id")

	// Get a product with product id
	// If there is no record, then throw a NotFound error
	var product models.Product
	if err := h.db.Take(&product, "id = ?", productID).Error; err != nil {
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
	product.Published = req.Published

	// Update product with data
	// If the update is unsuccessful, then throw an error
	if res := h.db.Save(&product); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, product)
}

/*
Description:

	Upload products image files to storage and update the product with the product id.

HTTP Method:

	POST `/api/v1/admin/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminProductHandler) UploadImages(c echo.Context) error {
	// Get product id from request
	productId := c.Param("id")

	// Get a product with product id
	// If there is no record, then throw a NotFound error
	var product models.Product
	if err := h.db.Take(&product, "id = ?", productId).Error; err != nil {
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
		if res := h.db.Create(&product_image); res.Error != nil {
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

	DELETE `/api/v1/admin/products/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminProductHandler) DeleteProduct(c echo.Context) error {
	// Get product id from request
	productID := c.Param("id")

	// Get a product with product id
	// If there is no record, then throw a NotFound error
	var product models.Product
	if err := h.db.Preload("ProductImages").Find(&product, "id = ?", productID).Error; err != nil {
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
			if res := h.db.Delete(&image); res.Error != nil {
				c.JSON(http.StatusInternalServerError, res.Error)
				return nil
			}
		}
	}

	// Delete a product
	// If the delete is unsuccessful, then throw an error
	if res := h.db.Delete(&product, "id = ?", productID); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the product")
}

/*
Description:

	Delete a specific product image with the product_image id and delete the object in storage.

HTTP Method:

	DELETE `/api/v1/admin/products/:id/assets/:image_id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminProductHandler) DeleteProductImage(c echo.Context) error {
	// Get product id from request
	productID := c.Param("id")

	// Get product image id from request
	productImageID := c.Param("image_id")

	// Get a product image with product id and product image id
	// If there is no record, then throw a NotFound error
	var productImage models.ProductImage
	if err := h.db.Find(&productImage, "id = ? AND product_id = ?", productImageID, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Delete the corresponsing object from S3
	if err := utils.Delete(productImage.Url); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}

	// Delete the product image record
	if res := h.db.Delete(&productImage); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the product image")
}
