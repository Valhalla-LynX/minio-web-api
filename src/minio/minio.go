package minio

import (
	"github.com/minio/minio-go"
	"log"
)

type MIOClient struct {
	c *minio.Client
}

func newMIOClient(client *minio.Client) *MIOClient {
	return &MIOClient{
		c: client,
	}
}

func InitMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *MIOClient {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatal(err)
	}
	return newMIOClient(minioClient)
}

func (mio *MIOClient) GetBucketList() []minio.BucketInfo {
	buckets, err := mio.c.ListBuckets()
	if err != nil {
		log.Fatal(err)
	}
	return buckets
}
