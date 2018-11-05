# Cache gin's middleware

[![Build Status](https://travis-ci.org/gin-contrib/cache.svg)](https://travis-ci.org/gin-contrib/cache)
[![codecov](https://codecov.io/gh/gin-contrib/cache/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/cache)](https://goreportcard.com/report/github.com/gin-contrib/cache)
[![GoDoc](https://godoc.org/github.com/gin-contrib/cache?status.svg)](https://godoc.org/github.com/gin-contrib/cache)

Gin middleware/handler to enable Cache.

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/yeqown/cache
```

Import it in your code:

```go
import "github.com/yeqown/cache"
```

### Canonical example:

See the examples:

* use with handler decorator, [example1](example/inmemory)
* use with redis store as a middleware **PageCache**, [example2](example/redis)
* use with **gin.Context.Get** also be a middleware, [example3](example/gincontext)

```go
package main

import (
	"fmt"
	"time"

	"github.com/yeqown/cache"
	"github.com/yeqown/cache/persistence"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := persistence.NewInMemoryStore(time.Second)
	
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	// Cached Page
	r.GET("/cache_ping", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	}))

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
```

### use with gin.Context

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
	// Cache is simplily call c.Set functions
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