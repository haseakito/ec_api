package routes

import (
	"os"

	"github.com/stripe/stripe-go/v76"
	"gorm.io/gorm"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"github.com/haseakito/ec_api/auth"
	"github.com/haseakito/ec_api/database"
	"github.com/haseakito/ec_api/handlers"
	"github.com/haseakito/ec_api/handlers/admin"
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

	// Initialize the Stripe client
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

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

	// Set up public APIs
	publicAPIs(r, db)

	// Set up admin APIs
	adminAPIs(r, db)

	return e
}

func publicAPIs(r *echo.Group, db *gorm.DB) {
	// Stores APIs Group
	s := r.Group("/stores")
	{
		// Initialize the new StoreController
		storeCtrl := handlers.NewStoreHandler(db)

		// Store APIs
		s.GET("", storeCtrl.GetStores)
		s.GET("/:id", storeCtrl.GetStore)

		// Product APIs for Stores
		s.GET("/:id/products", storeCtrl.GetProducts)

		// Order APIs for Stores
		s.POST("/:id/checkout", storeCtrl.CreateOrder)
		s.GET("/:id/orders", storeCtrl.GetOrders)
	}

	// Products APIs Group
	p := r.Group("/products")
	{
		// Initialize the new ProductHandler
		productCtrl := handlers.NewProductHandler(db)

		// Product APIs
		p.GET("/:id", productCtrl.GetProduct)

		// Review APIs
		p.GET("/:id/reviews", productCtrl.GetReviews)
		p.POST("/:id/reviews", productCtrl.CreateReview)
		p.DELETE("/:id/reviews/:review_id", productCtrl.DeleteReview)
	}

	// Webhooks Group
	w := r.Group("/webhooks")
	{
		// Initialize the new WebhookHandler
		webhookCtrl := NewWebhookHandler(db)

		w.POST("", webhookCtrl.StripeWebhook)
	}
}

func adminAPIs(r *echo.Group, db *gorm.DB) {
	// Set the admin API route
	a := r.Group("/admin")
	{
		/* Stores Group APIs */

		// Initialize the new AdminStoreHandler
		storeCtrl := admin.NewAdminStoreHandler(db)

		// Stores APIs
		a.POST("/stores", storeCtrl.CreateStore)
		a.PATCH("/stores/:id", storeCtrl.UpdateStore)
		a.DELETE("/stores/:id", storeCtrl.DeleteStore)

		// Assets APIs for Stores
		a.POST("/stores/:id/upload", storeCtrl.UploadImage)
		a.DELETE("/stores/:id/assets", storeCtrl.DeleteImage)

		// Product APIs for Stores
		a.POST("/stores/:id/products", storeCtrl.CreateProduct)
		a.GET("/stores/:id/products", storeCtrl.GetProducts)

		// Order APIs for Stores
		a.GET("/stores/:id/orders", storeCtrl.GetRevenues)

		/* Product Group APIs */

		// Initialize the new AdminProductHandler
		productCtrl := admin.NewAdminProductHandler(db)

		// Product APIs
		a.PATCH("/products/:id", productCtrl.UpdateProduct)
		a.POST("/products/:id/upload", productCtrl.UploadImages)
		a.DELETE("/products/:id", productCtrl.DeleteProduct)
		a.DELETE("/products/:id/assets/:image_id", productCtrl.DeleteProductImage)
	}
}
