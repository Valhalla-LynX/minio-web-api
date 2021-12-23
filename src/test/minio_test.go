package test

import (
	"github.com/minio/minio-go"
	"testing"
)

func TestMinio(t *testing.T) {
	endpoint := "192.168.50.222:29000"
	accessKeyID := "odm"
	secretAccessKey := "pwdodm2020"

	// 初使化 minio client对象。
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		t.Error(err)
	}

	buckets, err := minioClient.ListBuckets()
	if err != nil {
		t.Error(err)
		return
	}
	for _, bucket := range buckets {
		t.Log(bucket)
	}

	//log.Printf("%#v\n", minioClient) // minioClient初使化成功
}
