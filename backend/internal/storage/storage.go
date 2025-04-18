package storage

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sharithg/civet/internal/config"
)

type Storage struct {
	Client *minio.Client
}

func NewStorage(config *config.Config) *Storage {
	client := NewMinio(config)
	return &Storage{Client: client}
}

func (s *Storage) Upload(ctx context.Context, bucketName string, objectName string, filePath string, contentType string) (*minio.UploadInfo, error) {
	info, err := s.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *Storage) UploadImageBytes(ctx context.Context, bucketName, objectName string, imageBytes []byte, contentType string) (*minio.UploadInfo, error) {
	exists, err := s.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return nil, err
		}
	}

	reader := bytes.NewReader(imageBytes)
	size := int64(len(imageBytes))

	info, err := s.Client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *Storage) GetObjectBytes(ctx context.Context, bucketName string, objectName string, filePath string, contentType string) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "key-*.pem")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	err = s.Client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func (s *Storage) DeleteObject(ctx context.Context, bucketName string, objectName string) error {
	err := s.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListObjects(ctx context.Context, bucketName string) ([]string, error) {
	var objs []string
	opts := minio.ListObjectsOptions{}

	objectCh := s.Client.ListObjects(ctx, bucketName, opts)

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objs = append(objs, object.Key)
	}

	return objs, nil
}

func (s *Storage) CreateBucket(ctx context.Context, bucketName string) error {

	err := s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return err
	}

	log.Printf("Successfully created %s\n", bucketName)
	return nil
}

func (s *Storage) SetupBucket(ctx context.Context, bucketName string) error {
	exists, err := s.Client.BucketExists(ctx, bucketName)

	if err != nil {
		return err
	}

	if !exists {
		if err := s.CreateBucket(ctx, bucketName); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) SetupBuckets() error {
	ctx := context.Background()

	if err := s.SetupBucket(ctx, "node-pem-files"); err != nil {
		return err
	}

	if err := s.SetupBucket(ctx, "avatars"); err != nil {
		return err
	}

	return nil
}
func (s *Storage) GetObjectUrl(ctx context.Context, bucketName string, objectName string) (string, error) {
	url, err := s.Client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
