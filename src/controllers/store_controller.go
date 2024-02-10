package controllers

import (
	"errors"
	"net/http"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/haseakito/ec_api/storage"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StoreController struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new StoreController with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*StoreController: A pointer to the newly created StoreController instance.
*/
func NewStoreController(db *gorm.DB) *StoreController {
	return &StoreController{
		db: db,
	}
}

/*
Description:

	Creates a new store based on the data provided in the request payload.

HTTP Method:

	POST `/api/v1/stores`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) CreateStore(c echo.Context) error {
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
	if res := sc.db.Create(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusCreated, store)
}

/*
Description:

	Get all stores. Return empty array if no record is found.

HTTP Method:

	GET `/api/v1/stores`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) GetStores(c echo.Context) error {
	// Get all stores
	var stores []models.Store
	res := sc.db.Find(&stores)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, stores)
}

/*
Description:

	Get a specific store with the store id provided. Return nil if no record is found.

HTTP Method:

	GET `/api/v1/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) GetStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	var store models.Store
	res := sc.db.First(&store, "id = ?", storeId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Update a specific store with the store id provided and based on the data in the request payload.

HTTP Method:

	PATCH `/api/v1/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) UpdateStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	var store models.Store
	res := sc.db.First(&store, "id = ?", storeId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
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
	store.Name = req.Name
	store.Description = &req.Description

	// Update store with data
	// If the update is unsuccessful, then throw an error
	if res := sc.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Upload a store image file to storage and update the store with the store id.

HTTP Method:

	POST `/api/v1/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) UploadImage(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := sc.db.First(&store, "id = ?", storeId); err != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, nil)
			return nil
		}
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

	// Upload file to AWS S3 bucket
	// If the upload is unsuccessful, then throw an error
	url, err := storage.Upload(file, "stores/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}

	// Update the image url
	store.ImageUrl = &url

	// Update store with data
	// If the update is unsuccessful, then throw an error
	if res := sc.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Delete a specific store with the store id.

HTTP Method:

	DELETE `/api/v1/stores/:id`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (sc StoreController) DeleteStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	var store models.Store
	res := sc.db.Find(&store, "id = ?", storeId)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// If the store has an image url, delete the corresponding object from S3
	if store.ImageUrl != nil {
		storage.Delete(*store.ImageUrl)
	}

	// Delete a store
	// If the delete is unsuccessful, then throw an error
	if res := sc.db.Delete(&store, "id = ?", storeId); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the store")
}