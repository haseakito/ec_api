package requests

import validation "github.com/go-ozzo/ozzo-validation"

type StoreCreateRequest struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Description:

	Perform validation on the StoreCreateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r StoreCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.UserID,
			validation.Required.Error("User Id is required"),
		),
		validation.Field(
			&r.Name,
			validation.Length(0, 30),
			validation.Required.Error("Name is required"),
		),
		validation.Field(
			&r.Description,
			validation.Length(0, 1000),
		),
	)
}

type StoreUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Description:

	Perform validation on the StoreUpdateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r StoreUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Length(0, 30),
		),
		validation.Field(
			&r.Description,
			validation.Length(0, 1000),
		),
	)
}
