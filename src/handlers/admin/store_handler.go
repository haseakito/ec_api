package admin

import (
	"net/http"
	"time"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/haseakito/ec_api/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminStoreHandler struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new AdminStoreHandler with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*AdminStoreHandler: A pointer to the newly created AdminStoreHandler instance.
*/
func NewAdminStoreHandler(db *gorm.DB) *AdminStoreHandler {
	return &AdminStoreHandler{
		db: db,
	}
}

/*
Description:

	Creates a new store based on the data provided in the request payload.

HTTP Method:

	POST `/api/v1/admin/stores`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) CreateStore(c echo.Context) error {
	// Parsing request payload and validate the data
	// If there is a problem with the request, throw an error
	var req requests.StoreCreateRequest
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

	// Instantiate a new store
	store := models.Store{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: &req.Description,
	}

	// Create a new store
	// If the creation is unsuccessful, then throw an error
	if res := h.db.Create(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusCreated, store)
}

/*
Description:

	Upload a store image file to storage and update the store with the store id.

HTTP Method:

	POST `/api/v1/admin/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) UploadImage(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Get file from request
	// If there is a problem with the request, throw an error
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Validate request file
	// If there is a problem with the request, throw an error
	if err := requests.ValidateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	// If the store has an image url, delete the corresponding object from S3
	if store.ImageUrl != nil {
		utils.Delete(*store.ImageUrl)
	}

	// Upload file to AWS S3 bucket
	// If the upload is unsuccessful, then throw an error
	url, err := utils.Upload(file, "stores/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}

	// Update the image url
	store.ImageUrl = &url

	// Update store with data
	// If the update is unsuccessful, then throw an error
	if res := h.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Update a specific store with the store id provided and based on the data in the request payload.

HTTP Method:

	PATCH `/api/v1/admin/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) UpdateStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeId).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Parsing request payload and validate the data
	// If there is a problem with the request, throw an error
	var req requests.StoreUpdateRequest
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

	// Update store fields
	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Description != "" {
		store.Description = &req.Description
	}

	// Update store with data
	// If the update is unsuccessful, then throw an error
	if res := h.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Delete a specific store with the store id and delete the corresponding object in storage.

HTTP Method:

	DELETE `/api/v1/admin/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) DeleteStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeId).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// If the store has an image url, delete the corresponding object from S3
	if store.ImageUrl != nil {
		utils.Delete(*store.ImageUrl)
	}

	// Delete a store
	// If the delete is unsuccessful, then throw an error
	if res := h.db.Delete(&store, "id = ?", storeId); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the store")
}

/*
Description:

	Delete a store image with the store id.

HTTP Method:

	DELETE `/api/v1/admin/stores/:id/assets`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) DeleteImage(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeId).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// If the store has an image url, delete the corresponding object from S3
	if store.ImageUrl != nil {
		utils.Delete(*store.ImageUrl)
	}

	store.ImageUrl = nil

	if res := h.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully removed the store image")
}

/*
Description:

	Create a product for a specific store with the store id and based on the data provided in the request payload.

HTTP Method:

	GET `/api/v1/admin/products/:storeID`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) CreateProduct(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Parsing request payload and validate the data
	// If there is a problem with the request, throw an error
	var req requests.ProductCreateRequest
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

	// Instantiate a new product
	product := models.Product{
		StoreID:     store.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Create a new product for the store
	if res := h.db.Create(&product); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusCreated, product)
}

/*
Description:

	Get all products for a specific store with the store id.

HTTP Method:

	GET `/api/v1/admin/stores/:id/products`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h AdminStoreHandler) GetProducts(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	//
	// offsetStr := c.QueryParam("offset")
	// offset, err := strconv.Atoi(offsetStr)
	// if err != nil {
	// 	return echo.ErrBadRequest
	// }

	// Get a store with store id and products associated with the store
	// If there is no record, then throw a NotFound error
	var products []models.Product
	if err := h.db.Preload("ProductImages").Where("store_id = ?", storeID).Limit(10).Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, products)
}

func (h AdminStoreHandler) GetRevenues(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	var orders []models.Order
	if err := h.db.Preload("OrderItems.Product").Where("store_id = ? AND paid = ? AND created_at >= ?", storeID, true, oneYearAgo).Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	var totalRevenue float64
	for _, order := range orders {
		for _, item := range order.OrderItems {
			totalRevenue += float64(*item.Product.Price)
		}
	}

	res := map[string]interface{}{
		"orders":        orders,
		"total_revenue": totalRevenue,
		"sales_count":   len(orders),
	}

	return c.JSON(http.StatusOK, res)
}
