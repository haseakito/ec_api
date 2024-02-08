package routes

import (
	"github.com/haseakito/ec_api/controllers"
	"github.com/haseakito/ec_api/database"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	// Initialize database client
	db := database.Init()

	// Initialize new Echo application
	e := echo.New()

	// Set the default API route
	r := e.Group("/api/v1")

	// Stores APIs Group
	s := r.Group("/stores")
	{
		// Initialize the new StoreController
		storeCtrl := controllers.NewStoreController(db)

		// Define HTTP method to corresponding controller
		s.GET("", storeCtrl.GetStores)
		s.GET("/:id", storeCtrl.GetStore)
		s.POST("", storeCtrl.CreateStore)
		s.POST("/:id", storeCtrl.UploadImage)
		s.PATCH("/:id", storeCtrl.UpdateStore)
		s.DELETE("/:id", storeCtrl.DeleteStore)
	}

	return e
}
