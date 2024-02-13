package routes

import (
	"github.com/haseakito/ec_api/controllers"
	"github.com/haseakito/ec_api/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
Description:

	Initialize the Echo application, set up CORS configs, middlewares, routes for various APIs endpoints.

Returns:

	*echo.Echo: A pointer to the initialized Echo application instance.
*/
func Init() *echo.Echo {
	/* Initialize clients */

	// Initialize database client
	db := database.Init()

	// Initialize sentry client
	// TODO: Add Sentry SDK

	// Initialize clerk client
	// TODO: Add Clerk SDK

	// Initialize new Echo application
	e := echo.New()

	/* Configure middlewares */

	// Logger middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())

	// Recover from panics middleware
	e.Use(middleware.Recover())


	// Auth middleware

	/* Configure routers */

	// Set the default API route
	r := e.Group("/api/v1")

	// Stores APIs Group
	s := r.Group("/stores")
	{
		// Initialize the new StoreController
		storeCtrl := controllers.NewStoreController(db)

		// Define HTTP method to corresponding StoreController
		s.GET("", storeCtrl.GetStores)
		s.GET("/:id", storeCtrl.GetStore)
		s.GET("/:id/products", storeCtrl.GetProducts)
		s.POST("", storeCtrl.CreateStore)
		s.POST("/:id/products", storeCtrl.CreateProduct)
		s.POST("/:id/upload", storeCtrl.UploadImage)
		s.PATCH("/:id", storeCtrl.UpdateStore)
		s.DELETE("/:id", storeCtrl.DeleteStore)
	}

	// Products APIs Group
	p := r.Group("/products")
	{
		// Initialize the new ProductController
		productCtrl := controllers.NewProductController(db)

		// Define HTTP method to corresponding ProductController
		p.GET("/:id", productCtrl.GetProduct)
		p.POST("/:id/upload", productCtrl.UploadImages)
		p.PATCH("/:id", productCtrl.UpdateProduct)
		p.DELETE("/:id", productCtrl.DeleteProduct)
		p.DELETE("/:id/assets/:image_id", productCtrl.DeleteProductImage)
	}

	return e
}
