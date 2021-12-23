package token

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

type User struct {
	UserName    string `form:"UserName" json:"UserName" xml:"UserName"  binding:"required"`
	UserKey     string `form:"UserKey" json:"UserKey" xml:"UserKey" binding:"required"`
	ServiceName string `form:"ServiceName" json:"ServiceName" xml:"ServiceName" binding:"required"`
}

var allowKey map[string]string

func NewTokenMiddleware(secretKey string, identityKey string, allowKeyString string) *jwt.GinJWTMiddleware {
	allowKey = make(map[string]string)
	allowKeyStrings := strings.Split(allowKeyString, ",")
	for _, key := range allowKeyStrings {
		allowKey[key] = key
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "token",
		Key:         []byte(secretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ServiceName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVal User
			if err := c.ShouldBindJSON(&loginVal); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			return loginVal, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//todo token auth
			if v, ok := data.(*User); ok && v.ServiceName == "odm" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal(err)
	}
	return authMiddleware
}
