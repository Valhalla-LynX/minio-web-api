package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"gopkg.in/ini.v1"

	token "ResourcesConnector/src/middleware/token"
	mio "ResourcesConnector/src/minio"

	"log"
	"os"
)

var mioClient *mio.MIOClient

func main() {
	wd, _ := os.Getwd()
	cfg, err := ini.Load(wd + "/config.ini")
	if err != nil {
		log.Fatal(err)
		return
	}

	endpoint := cfg.Section("minio").Key("endpoint").String()
	accessKeyID := cfg.Section("minio").Key("accessKeyID").String()
	secretAccessKey := cfg.Section("minio").Key("secretAccessKey").String()
	useSSL, err := cfg.Section("minio").Key("useSSL").Bool()
	if err != nil {
		log.Fatal("\"useSSL\" must be true or false")
		return
	}

	port := cfg.Section("gin").Key("port").String()

	secretKey := cfg.Section("token").Key("secretKey").String()
	identityKey := cfg.Section("token").Key("identityKey").String()
	allowKey := cfg.Section("key").Key("name").String()

	mioClient = mio.InitMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL)

	token := token.NewTokenMiddleware(secretKey, identityKey, allowKey)

	r := gin.Default()

	errInit := token.MiddlewareInit()
	if errInit != nil {
		log.Fatal("Init token middleware failed:" + errInit.Error())
	}

	r.NoRoute(token.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/login", token.LoginHandler)
	auth := r.Group("/auth")
	auth.GET("/refresh_token", token.RefreshHandler)
	auth.Use(token.MiddlewareFunc())
	{
		auth.GET("/getBucketList", getBucketList)
		auth.GET("/getObject", getObject)
	}

	err = r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

func getBucketList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": mioClient.GetBucketList(),
	})
}
func getObjectList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": mioClient.GetBucketList(),
	})
}
func getObject(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": mioClient.GetBucketList(),
	})
}
func getObjectByBase64(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": mioClient.GetBucketList(),
	})
}
func putObject(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": mioClient.GetBucketList(),
	})
}
