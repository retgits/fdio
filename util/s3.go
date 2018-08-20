// Package util implements utility methods
package util

// The imports
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloadFile downloads a file from Amazon S3 and stores it  in the specified location.
func DownloadFile(awsSession *session.Session, folder string, filename string, bucket string) error {
	// Create an instance of the S3 Downloader
	s3Downloader := s3manager.NewDownloader(awsSession)

	// Create a new temporary file
	tempFile, err := os.Create(filepath.Join(folder, filename))
	if err != nil {
		return err
	}

	// Prepare the download
	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}

	// Download the file to disk
	_, err = s3Downloader.Download(tempFile, objectInput)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile creates a copy of an existing file with a new name
func CopyFile(awsSession *session.Session, filename string, bucket string) error {
	// Create an instance of the S3 Session
	s3Session := s3.New(awsSession)

	// Prepare the copy object
	objectInput := &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", bucket, filename)),
		Key:        aws.String(fmt.Sprintf("%s_bak", filename)),
	}

	// Copy the object
	_, err := s3Session.CopyObject(objectInput)
	if err != nil {
		return err
	}

	return nil
}

// UploadFile uploads a file to Amazon S3
func UploadFile(awsSession *session.Session, folder string, filename string, bucket string) error {
	// Create an instance of the S3 Uploader
	s3Uploader := s3manager.NewUploader(awsSession)

	// Create a file pointer to the source
	reader, err := os.Open(filepath.Join(folder, filename))
	if err != nil {
		return err
	}
	defer reader.Close()

	// Prepare the upload
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   reader,
	}

	// Upload the file
	_, err = s3Uploader.Upload(uploadInput)
	if err != nil {
		return err
	}

	return nil
}
