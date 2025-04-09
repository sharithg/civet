package storage

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sharithg/civet/internal/config"
)

func NewMinio(config *config.Config) *minio.Client {
	endpoint := config.MinioHost
	accessKeyID := config.MinioAccessKey
	secretAccessKey := config.MinioSecretKey
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal("error configuring minio: ", err)
	}

	return minioClient
}
