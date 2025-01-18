package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Uploads a file to a folder in S3 and returns its public URL
func UploadToS3(file *multipart.FileHeader, folder string) (string, error) {
	bucketName := "eliotgoldenarm"
	region := "us-east-2"
	key := fmt.Sprintf("%s/%s", folder, file.Filename)

	// Open the file
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the file into a buffer to detect the content type
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(fileBytes)

	// Reopen the file for uploading
	f, err = file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create a new AWS session (automatically picks up environment variables or IAM roles)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// Create an S3 uploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3 with the correct headers
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(key),
		Body:               f,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("inline"),
	})
	if err != nil {
		return "", err
	}

	// Return the public URL
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key), nil
}
