package utils

import (
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/*
Description:

	Upload a file to AWS S3 bucket and returns the URL of the uploaded file

Parameters:

	file (*multipart.FileHeader): The file to be uploaded
	key (string): The key prefix for the file in S3.

Returns:

	(string, error): The URL of the uploaded file. Otherwise, any error encountered during the upload process.
*/
func Upload(file *multipart.FileHeader, key string) (string, error) {
	// Initialize AWS session
	// If initialize failed, then throw an error
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to create AWS session: ", err)
		return "", err
	}

	// Initialize S3 upload client
	uploader := s3manager.NewUploader(sess)

	// Open file
	// If opening file is unsuccessful, throw an error
	src, err := file.Open()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Upload file to AWS S3
	// If there a problem with uploading, throw an error
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key + url.QueryEscape(file.Filename)),
		Body:   src,
	})
	if err != nil {
		log.Fatal("failed to upload object to S3: ", err)
		return "", err
	}

	// Return the image url location
	return res.Location, nil
}

/*
Description:

	Delete a file in AWS S3 bucket.

Parameters:

	imageUrl (string): The URL of the file to be deleted.

Returns:

	error: Any error encountered during the upload process.
*/
func Delete(imageUrl string) error {
	// Parse the S3 object key from the image URL
	key := imageUrl[strings.Index(imageUrl, "amazonaws.com/")+len("amazonaws.com/"):]

	// Initialize AWS session
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to create AWS session: ", err)
		return err
	}

	// Initialize S3 client
	s3Client := s3.New(sess)

	// Delete the object from S3
	// If there is a problem with deleting, throw an error
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal("failed to delete object from S3: ", err)
		return err
	}

	return nil
}
