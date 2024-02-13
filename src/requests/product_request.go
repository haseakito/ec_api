package requests

import validation "github.com/go-ozzo/ozzo-validation"

type ProductCreateRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Price       *float32 `json:"price"`
}

/*
Description:

	Perform validation on the ProductCreateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r ProductCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Product name is requried"),
			validation.Length(0, 255),
		),
		validation.Field(
			&r.Description,
			validation.Length(0, 1000),
		),
		validation.Field(
			&r.Price,
		),
	)
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Published   bool    `json:"is_published"`
}

/*
Description:

	Perform validation on the ProductUpdateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r ProductUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Length(0, 255),
		),
		validation.Field(
			&r.Description,
			validation.Length(0, 1000),
		),
		validation.Field(
			&r.Price,
		),
		validation.Field(
			&r.Published,
		),
	)
}
