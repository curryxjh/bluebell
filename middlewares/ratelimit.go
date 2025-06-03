package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit ...")
			c.Abort()
			return
		}
		c.Next()
	}
}
