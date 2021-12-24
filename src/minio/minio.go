package minio

import (
	m "github.com/minio/minio-go"
	"log"
)

func InitMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *m.Client {
	minioClient, err := m.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatal(err)
	}
	return minioClient
}
