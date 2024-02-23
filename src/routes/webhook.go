package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"gorm.io/gorm"

	"github.com/haseakito/ec_api/models"
)

type WebhookHandler struct {
	db *gorm.DB
}

/*
Description:

	Instantiates a new WebhookHandler with the provided database connection.

Parameters:

	db (*gorm.DB): A pointer to the GORM database connection.

Returns:

	*WebhookHandler: A pointer to the newly created WebhookHandler instance.
*/
func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{
		db: db,
	}
}

/*
Description:

	StripeWebhook handles incoming webhook events from Stripe.
	It processes the checkout.session.completed event to update the order status.

HTTP Method:

	POST `/api/v1/webhooks`

Parameters:

	c (echo.Context): Context object containing the HTTP request information.

Returns:

	An error if any occurred during the execution of the function, nil otherwise.
*/
func (h WebhookHandler) StripeWebhook(c echo.Context) error {
	// Set the maximum body bytes for the request
	const MaxBodyBytes = int64(65536)
	//
	body := http.MaxBytesReader(c.Response().Writer, c.Request().Body, MaxBodyBytes)

	// Read the request body
	payload, err := io.ReadAll(body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to read webhook body")
	}

	// Get the Stripe signature from the request header
	sigHeader := c.Request().Header.Get("Stripe-Signature")

	// Construct the Stripe event from the payload
	event, err := webhook.ConstructEvent(payload, sigHeader, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to construct stripe event")
	}

	// Check if the event type is checkout.session.completed
	if event.Type == "checkout.session.completed" {
		// Unmarshal the event data into a PaymentIntent object
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// Extract the order ID from the PaymentIntent's metadata
		orderID := paymentIntent.Metadata["order_id"]

		// Get an order with order id
		// If there is no record, then throw a NotFound error
		var order models.Order
		if err := h.db.Find(&order, "id = ?", orderID).Error; err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		// Update order field
		order.Paid = true

		// Update order with data
		// If the update is unsuccessful, then throw an error
		if err := h.db.Save(&order).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, "Successfully updated the order")
	}
	return nil
}
