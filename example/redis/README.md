### redis usage

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeqown/cache"
	"github.com/yeqown/cache/persistence"
)

func main() {
	r := gin.Default()

	store := persistence.NewRedisCache("127.0.0.1:6379", "", 1, time.Second*30)
	// no cache page
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Cached Page
	r.GET("/cache_ping", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	}))

	r.GET("/cache_ping_json", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"timestamp": fmt.Sprint(time.Now().Unix())})
	}))

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

```