package main

import (
	mio "ResourcesConnector/src/minio"
	"github.com/ReneKroon/ttlcache/v2"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"gopkg.in/ini.v1"

	"ResourcesConnector/src/middleware/cache"
	"ResourcesConnector/src/middleware/token"

	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var endpoint string
var useSSL bool
var mioClient *minio.Client
var mioCacheKey ttlcache.SimpleCache

func main() {
	wd, _ := os.Getwd()
	cfg, err := ini.Load(wd + "/config.ini")
	if err != nil {
		log.Fatal(err)
		return
	}

	endpoint = cfg.Section("minio").Key("endpoint").String()
	accessKeyID := cfg.Section("minio").Key("accessKeyID").String()
	secretAccessKey := cfg.Section("minio").Key("secretAccessKey").String()
	useSSL, err = cfg.Section("minio").Key("useSSL").Bool()
	if err != nil {
		log.Fatal("\"useSSL\" must be true or false")
		return
	}

	port := cfg.Section("gin").Key("port").String()

	secretKey := cfg.Section("tokenMiddleware").Key("secretKey").String()
	identityKey := cfg.Section("tokenMiddleware").Key("identityKey").String()
	allowKey := cfg.Section("key").Key("name").String()

	mioClient = mio.InitMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	mioCacheKey = ttlcache.NewCache()
	mioCacheKey.SetTTL(24 * 60 * time.Minute)

	tokenMiddleware := token.NewTokenMiddleware(secretKey, identityKey, allowKey)
	cacheMiddleware30m := cache.NewTokenMiddleware(30 * time.Minute)
	//cacheMiddleware10m := cache.NewTokenMiddleware(10 * time.Minute)
	//cacheMiddleware3m := cache.NewTokenMiddleware(3 * time.Minute)
	cacheMiddleware1m := cache.NewTokenMiddleware(1 * time.Minute)

	r := gin.Default()

	errInit := tokenMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("Init tokenMiddleware middleware failed:" + errInit.Error())
	}

	r.NoRoute(tokenMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/login", tokenMiddleware.LoginHandler)
	auth := r.Group("/auth")
	auth.GET("/refresh_token", tokenMiddleware.RefreshHandler)
	auth.Use(tokenMiddleware.MiddlewareFunc())
	{
		auth.GET("/getBucketList", cacheMiddleware30m, getBucketList)
		auth.GET("/getObject", cacheMiddleware1m, getObject)
		auth.POST("/postObject", postObject)
	}

	err = r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

func getBucketList(c *gin.Context) {
	buckets, err := mioClient.ListBuckets()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"message": buckets,
	})
}

func postObject(c *gin.Context) {
	bucketName := c.Query("bucket")

	file, _ := c.FormFile("file")
	fileReader, _ := file.Open()

	n, err := mioClient.PutObject(bucketName, file.Filename, fileReader, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"message": n,
	})
}

func getObject(c *gin.Context) {
	bucketName := c.Query("bucket")
	objectName := c.Query("name")
	//"http://192.168.50.222:29000/odm1/star.jpeg?"
	var buffer bytes.Buffer
	if useSSL {
		buffer.WriteString("https://")
	} else {
		buffer.WriteString("http://")
	}
	buffer.WriteString(endpoint)
	buffer.WriteString("/")
	buffer.WriteString(bucketName)
	buffer.WriteString("/")
	buffer.WriteString(objectName)
	buffer.WriteString("?")
	reqHost := buffer.String()
	// get cacheKetName or sign a new one
	ketValue, err := getObjectCacheValue(reqHost, bucketName, objectName)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	// ask object
	response, err := http.Get(reqHost + ketValue)
	// local cache expired or other error
	if err == nil && response.StatusCode == http.StatusForbidden {
		// resign token
		ketValue, err = signObjectCacheValue(reqHost, bucketName, objectName)
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		response, err = http.Get(reqHost + ketValue)
	}

	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	defer reader.Close()
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, nil)
}

func getObjectCacheValue(objectCacheKey string, bucketName string, objectName string) (string, error) {
	if val, err := mioCacheKey.Get(objectCacheKey); err != ttlcache.ErrNotFound {
		return fmt.Sprintf("%s", val), err
	} else {
		return signObjectCacheValue(objectCacheKey, bucketName, objectName)
	}
}

func signObjectCacheValue(objectCacheKey string, bucketName string, objectName string) (string, error) {
	sign, err := mioClient.PresignedGetObject(bucketName, objectName, 1.5*24*60*time.Minute, nil)
	mioCacheKey.Set(objectCacheKey, sign.RawQuery)
	return sign.RawQuery, err
}
