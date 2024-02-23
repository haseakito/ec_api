package requests

import validation "github.com/go-ozzo/ozzo-validation"

type CheckoutCreateRequest struct {
	UserID    string   `json:"user_id"`
	ProductIDs []string `json:"product_ids"`
}

/*
Description:

	Perform validation on the CheckoutCreateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r CheckoutCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.UserID,
			validation.Required.Error("User Id is required"),
		),
		validation.Field(
			&r.ProductIDs,
			validation.Required.Error("Product Ids is required"),
		),
	)
}
