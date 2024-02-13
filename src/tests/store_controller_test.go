package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/haseakito/ec_api/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateStore(t *testing.T) {
	// Initialize new Echo application
	e := routes.Init()

	// JSON data to create a store
	storeJSON := `{"name": "test store", "description": "test store description"}`

	// Create a new HTTP request with the store json data
	req := httptest.NewRequest(http.MethodPost, "/api/v1/stores", strings.NewReader(storeJSON))
	// Set the Content-Type in header to be application/json
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// POST request to the specified endpoint
	e.ServeHTTP(rec, req)

	//
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestGetStores(t *testing.T) {
	// Initialize new Echo application
	e := routes.Init()

	// Create a new HTTP request with the store json data
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stores", nil)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// GET request to the specified endpoint
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// TODO: modify the store id
func TestGetStore(t *testing.T) {
	// Initialize new Echo application
	e := routes.Init()

	// Create a new HTTP request with the store json data
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stores/store-123", nil)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// GET request to the specified endpoint
	e.ServeHTTP(rec, req)
}

// TODO: modify the store id
func TestUpdateStore(t *testing.T) {
	// Initialize new Echo application
	e := routes.Init()

	// JSON data to create a store
	storeJSON := `{"name": "updated test store", "description": "updated test store description"}`

	// Create a new HTTP request with the store json data
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/stores/store-123", strings.NewReader(storeJSON))
	// Set the Content-Type in header to be application/json
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// PATCH request to the specified endpoint
	e.ServeHTTP(rec, req)

	//
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestDeleteStore(t *testing.T) {
	// Initialize new Echo application
	e := routes.Init()

	// Create a new HTTP request with the store json data
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/stores/store-123", nil)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// PATCH request to the specified endpoint
	e.ServeHTTP(rec, req)

	//
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equalf(t, "Successfully deleted the store", rec.Body.String(), "error message %s")
}
