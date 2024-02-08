package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StoreController struct {
	db *gorm.DB
}

func NewStoreController(db *gorm.DB) *StoreController {
	return &StoreController{
		db: db,
	}
}

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

func (sc StoreController) UploadImage(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

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

	// Initialize AWS session
	// If initialize failed, then throw an error
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Failed to create session:", err)
		return nil
	}

	// Initialize S3 upload client
	uploader := s3manager.NewUploader(sess)

	// Open file
	// If opening file is unsuccessful, 
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Upload file to AWS S3
	// If there a problem with uploading, throw an error
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String("stores/" + file.Filename),
		Body:   src,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := sc.db.First(&store, "id = ?", storeId); err != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, nil)
			return nil
		}
	}

	// Update the image url
	store.ImageUrl = &res.Location

	// Update store with data
	// If the update is unsuccessful, then throw an error
	if res := sc.db.Save(&store); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

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

	// If the store has an image URL, delete the corresponding object from S3
	if store.ImageUrl != nil {
		// Parse the S3 object key from the image URL
		imageURL := *store.ImageUrl
		key := imageURL[strings.Index(imageURL, "amazonaws.com/")+len("amazonaws.com/"):]

		// Initialize AWS session
		sess, err := session.NewSession()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}
		
		// Initialize S3 client
		s3Client := s3.New(sess)

		// Delete the object from S3
		// If there is a problem with deleting, throw an error
		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:    aws.String(key),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return nil
		}
	}

	// Delete a store
	// If the delete is unsuccessful, then throw an error
	if res := sc.db.Delete(&store, "id = ?", storeId); res.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Error)
		return nil
	}

	return c.JSON(http.StatusOK, "Successfully deleted the store")
}
