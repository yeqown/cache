### use as custom Cache

```go
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeqown/cache"
	"github.com/yeqown/cache/persistence"
)

func main() {
	r := gin.Default()

	// link redis connection
	store := persistence.NewRedisCache("127.0.0.1:6379", "", 1, time.Second*30)

	cus := r.Group("/cus")
	// to set gin.Context With store(data saved in redis)
	cus.Use(cache.Cache(store))
	{
		cus.GET("/ex1", WithContextExampleHdl)
	}

	r.Run(":8080")
}

// WithContextExampleHdl self custom cache with cache
func WithContextExampleHdl(c *gin.Context) {
	// get cache firstly
	v, _ := c.Get(cache.CACHE_MIDDLEWARE_KEY)
	cache := v.(persistence.CacheStore)
	var timePtr = new(int64)
	// get from cache
	if err := cache.Get("time", timePtr); err != nil {

		// missed by the key, reset into cache
		if err == persistence.ErrCacheMiss {
			unix := time.Now().Unix()
			cache.Set("time", unix, time.Second*10)
			c.JSON(http.StatusOK, gin.H{"time": unix})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"time": *timePtr})
}

```