package s3

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client     *s3.Client
	BucketName string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		Client:     client,
		BucketName: bucket,
	}
}

func (s *S3Storage) UploadFile(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	key string,
) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
		Body:   file,
		ContentType: func() *string {
			v := "application/pdf"
			return &v
		}(),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://%s.s3.amazonaws.com/%s",
		s.BucketName,
		key,
	)

	return url, nil
}