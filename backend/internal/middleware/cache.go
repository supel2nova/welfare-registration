package middleware

import "github.com/gin-gonic/gin"

const (
	CacheRefData = "public, max-age=86400"
	CacheSearch  = "private, max-age=60"
)

func CacheControl(value string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", value)
		c.Next()
	}
}
