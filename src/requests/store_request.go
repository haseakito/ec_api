package requests

import validation "github.com/go-ozzo/ozzo-validation"

type StoreCreateRequest struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r StoreCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.UserID,
			validation.Required.Error("User Id is required"),
		),
		validation.Field(
			&r.Name,
			validation.Length(0, 255),
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