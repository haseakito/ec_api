package routes

import (
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"github.com/haseakito/ec_api/auth"
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
	client, _ := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))

	// Initialize new Echo application
	e := echo.New()

	/* Configure middlewares */

	// CORS middleware
	e.Use(middleware.CORS())

	// Logger middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())

	// Recover from panics middleware
	e.Use(middleware.Recover())

	// Auth middleware
	e.Use(auth.AuthMiddleware(client))

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
		s.DELETE("/:id/assets", storeCtrl.DeleteImage)
		s.DELETE("/:id", storeCtrl.DeleteStore)
	}

	// Products APIs Group
	p := r.Group("/products")
	{
		// Initialize the new ProductController
		productCtrl := controllers.NewProductController(db)

		// Define HTTP method to corresponding ProductController
		p.GET("/:id", productCtrl.GetProduct)
		p.GET("/:id/reviews", productCtrl.GetReviews)
		p.POST("/:id/reviews", productCtrl.CreateReview)
		p.POST("/:id/upload", productCtrl.UploadImages)
		p.PATCH("/:id", productCtrl.UpdateProduct)
		p.DELETE("/:id", productCtrl.DeleteProduct)
		p.DELETE("/:id/assets/:image_id", productCtrl.DeleteProductImage)
		p.DELETE("/:id/reviews/:review_id", productCtrl.DeleteReview)
	}

	return e
}
