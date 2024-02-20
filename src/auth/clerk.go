package auth

import (
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"github.com/labstack/echo/v4"
)

/*
Description:

	AuthMiddleware is used to authenticate incoming requests by verifying the session claims from the context.
	It retrieves the session claims from the context, validates them, and retrieves the user information from Clerk.
	If the session claims are valid and the user information is successfully retrieved, the user information is set in the context,
	and the request is passed to the next handler in the middleware chain.

Parameters:

	client (clerk.Client): The Clerk client used to interact with the Clerk authentication service.

Returns:

	echo.,MiddlewareFunc: An Echo middleware function that performs authentication for incoming requests.
*/
func AuthMiddleware(client clerk.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Obtain bearer token from request header
			sessToken := c.Request().Header.Get("Authorization")
			sessToken = strings.TrimPrefix(sessToken, "Bearer ")

			// Verify the bearer token
			// If the verication is unsuccessful, then throw an unauthroized error
			sessClaims, err := client.VerifyToken(sessToken)
			if err != nil {
				return echo.ErrUnauthorized
			}

			// Retrieve user information from Clerk
			user, err := client.Users().Read(sessClaims.Claims.Subject)
			if err != nil {
				return echo.ErrUnauthorized
			}

			// Set user information in the context
			c.Set("user", user)

			// Call the next handler in the middleware chain
			return next(c)
		}
	}
}
