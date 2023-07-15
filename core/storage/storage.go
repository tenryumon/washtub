package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	AWSS3BucketUserContent = "user-content"
)

type Storage struct {
	s3Client   *s3.Client
	bucketName *string
	domain     string
}

type Config struct {
	Region string
	Domain string
}

type Object struct {
	Key     string
	Content []byte
	Type    string
}

func New(storageConf Config) (*Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(storageConf.Region))
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client
	return &Storage{
		s3Client: s3.NewFromConfig(cfg),
		domain:   storageConf.Domain,
	}, nil
}

func (s *Storage) SetBucketName(ctx context.Context, bucket string) error {
	if _, err := s.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return err
	}
	s.bucketName = aws.String(bucket)
	return nil
}

func (s *Storage) UploadFile(ctx context.Context, object Object) error {
	buff := bytes.NewReader(object.Content)

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(AWSS3BucketUserContent),
		Key:         aws.String(object.Key),
		Body:        buff,
		ContentType: aws.String(object.Type),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GenerateUserContentURL(ctx context.Context, relativePath string) string {
	if relativePath == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.domain, relativePath)
}
