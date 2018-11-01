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

	store := persistence.NewInMemoryStore(60 * time.Second)
	// Cached Page
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	r.GET("/cache_ping", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	}))

	// Use as middleware
	cached := r.Group("/cached")
	cached.Use(cache.SiteCache(store, time.Minute))
	{
		cached.GET("/ping", func(c *gin.Context) {
			timestamp := fmt.Sprint(time.Now().Unix())
			c.JSON(http.StatusOK, gin.H{"timestamp": timestamp})
		})

		cached.GET("/ping_with_param", func(c *gin.Context) {
			type pingForm struct {
				Who string `form:"who" binding:"required"`
			}
			var (
				f         pingForm
				timestamp = fmt.Sprint(time.Now().Unix())
			)

			if err := c.ShouldBind(&f); err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error() + timestamp})
				return
			}

			c.JSON(http.StatusOK, gin.H{"timestamp": timestamp})
		})
	}

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
