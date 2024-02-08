package requests

import (
	"errors"
	"mime/multipart"
)

func ValidateFile(file *multipart.FileHeader) error {
	// If there is no file, throw an error
	if file == nil {
		return errors.New("File is required")
	}

	// Check file size (Max size 2MB limit)
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return errors.New("File size exceeds the limit of 5MB")
	}

	// Check file type (jpeg or png)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	// 
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return errors.New("file type must be JPEG or PNG")
	}

	return nil
}