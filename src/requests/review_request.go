package requests

import validation "github.com/go-ozzo/ozzo-validation"

type ReviewCreateRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

/*
Description:

	Perform validation on the ReviewCreateRequest struct fields.

Returns:

	error: An error if any validation fails, otherwise nil.
*/
func (r ReviewCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.UserID,
			validation.Required.Error("User Id is required"),
		),
		validation.Field(
			&r.Content,
			validation.Required.Error("Review content is required"),
			validation.Length(0, 255),
		),
	)
}
