package cache

import (
	c "github.com/Valhalla-LynX/gin-cache"
	"github.com/Valhalla-LynX/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
)

func NewTokenMiddleware(expire time.Duration) gin.HandlerFunc {
	memoryStore := persist.NewMemoryStore(time.Minute)
	return c.CacheByRequestURI(memoryStore, expire)
}
