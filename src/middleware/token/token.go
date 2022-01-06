package token

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"strconv"

	"errors"
	"log"
	"os"
	"strings"
	"time"
)

type User struct {
	UserName    string `form:"UserName" json:"UserName" xml:"UserName"  binding:"required"`
	UserKey     string `form:"UserKey" json:"UserKey" xml:"UserKey" binding:"required"`
	ServiceName string `form:"ServiceName" json:"ServiceName" xml:"ServiceName" binding:"required"`
	Bucket      string
	Permissions int
}

var allowUser map[string]User
var permissionURLs map[string]int
var secretKey string
var identityKey string

func initAuth() error {
	wd, _ := os.Getwd()
	cfg, err := ini.Load(wd + "/config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	secretKey = cfg.Section("token").Key("secretKey").String()
	identityKey = cfg.Section("token").Key("identityKey").String()
	name := cfg.Section("key").Key("name").String()
	pwd := cfg.Section("key").Key("pwd").String()

	permissionURLs = make(map[string]int)
	allowUser = make(map[string]User)

	URLCfg := cfg.Section("token").Key("permissionURLs").String()
	URLCfgArray := strings.Split(URLCfg, ";")

	for _, cfg := range URLCfgArray {
		url := strings.Split(cfg, ",")
		if len(url) > 1 {
			permission, err := strconv.Atoi(url[1])
			if err != nil {
				return errors.New("illegal permissionURLs")
			}
			permissionURLs[url[0]] = permission
		}
	}

	nameStrings := strings.Split(name, ",")
	pwdStrings := strings.Split(pwd, ",")

	if len(pwdStrings) < len(nameStrings) {
		return errors.New("illegal name&pwd")
	}

	for i, name := range nameStrings {
		user := User{
			UserName:    name,
			UserKey:     pwdStrings[i],
			Bucket:      cfg.Section(name).Key("bucket").String(),
			Permissions: cfg.Section(name).Key("permissions").MustInt(0),
		}
		allowUser[name] = user
	}
	return nil
}

func NewTokenMiddleware() (*jwt.GinJWTMiddleware, error) {
	err := initAuth()
	if err != nil {
		return nil, err
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "token",
		Key:         []byte(secretKey),
		Timeout:     7 * 24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(User); ok {
				return jwt.MapClaims{
					identityKey:   v.UserName,
					"ServiceName": v.ServiceName,
					"Bucket":      allowUser[v.UserName].Bucket,
					"Permissions": allowUser[v.UserName].Permissions,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName:    claims[identityKey].(string),
				ServiceName: claims["ServiceName"].(string),
				Bucket:      claims["Bucket"].(string),
				Permissions: (int)(claims["Permissions"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVal User
			if err := c.ShouldBindJSON(&loginVal); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if loginVal.UserKey != allowUser[loginVal.UserName].UserKey {
				return "", errors.New("wrong user key")
			}
			return loginVal, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			url := strings.Split(c.Request.URL.Path, "/")
			if permission, permissionOk := permissionURLs[url[len(url)-1]]; permissionOk {
				if v, ok := data.(*User); ok {
					if user, have := allowUser[v.UserName]; have {
						if (strings.Contains(user.Bucket, c.Query("bucket")) ||
							user.Bucket == "*") &&
							user.Permissions&permission == permission {
							return true
						}
					}
				}
				return false
			}
			return true
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
	return authMiddleware, nil
}
