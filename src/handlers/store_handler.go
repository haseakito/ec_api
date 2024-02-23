package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"gorm.io/gorm"

	"github.com/haseakito/ec_api/models"
	"github.com/haseakito/ec_api/requests"
)

type StoreHandler struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new StoreHandler with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*StoreHandler: A pointer to the newly created StoreHandler instance.
*/
func NewStoreHandler(db *gorm.DB) *StoreHandler {
	return &StoreHandler{
		db: db,
	}
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
func (h StoreHandler) GetStores(c echo.Context) error {
	// Get all stores
	// If there is no record, then throw a NotFound error
	var stores []models.Store
	if err := h.db.Find(&stores).Error; err != nil {
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
func (h StoreHandler) GetStore(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeId).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, store)
}

/*
Description:

	Get all published products for a specific store with the store id.

HTTP Method:

	GET `/api/v1/stores/:id/products`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h StoreHandler) GetProducts(c echo.Context) error {
	// Get store id from request
	storeId := c.Param("id")

	//
	// offsetStr := c.QueryParam("offset")
	// offset, err := strconv.Atoi(offsetStr)
	// if err != nil {
	// 	return echo.ErrBadRequest
	// }

	// Get a store with store id and products associated with the store
	var products []models.Product
	res := h.db.Preload("ProductImages").Where("store_id = ? AND published = ?", storeId, true).Order("created_at desc").Limit(10).Find(&products)

	// If there is no record, then throw a NotFound error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	return c.JSON(http.StatusOK, products)
}

/*
Description:

	Creates a new order based on the data provided in the request payload.

HTTP Method:

	POST `/api/v1/stores/:id/checkout`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, Stripe checkout session url otherwise.
*/
func (h StoreHandler) CreateOrder(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	// Get a store with store id
	// If there is no record, then throw a NotFound error
	var store models.Store
	if err := h.db.Take(&store, "id = ?", storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return nil
	}

	// Validate request data
	// If there is a problem with the request, throw an error
	var req requests.CheckoutCreateRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	// Iterate through product IDs to instantiate a new checkout session line items
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, productID := range req.ProductIDs {
		// Get a product with product id
		var product models.Product
		if err := h.db.Where("id = ? AND published = ?", productID, true).Take(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, err)
			return nil
		}

		// Instantiate a new checkout session line item
		lineItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(product.Name),
				},
				UnitAmount: stripe.Int64(int64(*product.Price * 100)),
			},
			Quantity: stripe.Int64(1),
		}

		// Add checkout session line item to array
		lineItems = append(lineItems, lineItem)
	}

	// Transaction to create an order and order items associated with the order
	// If the transaction failed, then throw an error
	var order models.Order
	h.db.Transaction(func(tx *gorm.DB) error {
		// Instantiate a new order
		order = models.Order{
			StoreID: storeID,
			UserID:  req.UserID,
			Paid:    false,
		}

		// Create a new order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Iterate through product IDs to create order items associated with the order
		for _, productID := range req.ProductIDs {
			// Instantiate a new order item
			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: productID,
			}

			// Create a new order item
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	// Instantiate a stripe checkout session
	params := &stripe.CheckoutSessionParams{
		LineItems: lineItems,
		Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
		Metadata: map[string]string{
			"order_id": order.ID,
		},
		SuccessURL: stripe.String(os.Getenv("FRONT_URL") + "/" + storeID + "/cart?success=true"),
		CancelURL:  stripe.String(os.Getenv("FRONT_URL") + "/" + storeID + "/cart?canceled=true"),
	}

	// Create a new stripe checkout session
	res, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}

	return c.JSON(http.StatusCreated, res.URL)
}

/*
Description:

	Get specific orders with the store id provided. Return nil if no record is found.

HTTP Method:

	GET `/api/v1/stores/:id/orders`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h StoreHandler) GetOrders(c echo.Context) error {
	// Get store id from request
	storeID := c.Param("id")

	// 
	var orders []models.Order
	if err := h.db.Preload("OrderItems.Product").Where("store_id = ?", storeID).Order("created_at desc").Limit(10).Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, orders)
}
