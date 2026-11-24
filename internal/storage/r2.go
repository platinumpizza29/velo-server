package storage

import (
	"context"
	"fmt"
	"time"
	"velo-server/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Storage struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

func NewR2Storage(config *model.Config) *R2Storage {
	cfg := aws.Config{
		Region: "auto",
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				config.R2AccessKeyID,
				config.R2SecretAccessKey,
				"",
			),
		),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(config.R2Endpoint)
		o.UsePathStyle = true
	})

	return &R2Storage{
		client:    client,
		presigner: s3.NewPresignClient(client),
		bucket:    config.R2BucketName,
	}
}

func (r *R2Storage) GeneratePresignedURL(
	ctx context.Context,
	objectKey string,
) (string, error) {

	result, err := r.presigner.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(15*time.Minute),
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return result.URL, nil
}
