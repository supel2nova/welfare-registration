package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				status, env := httpx.Fail(apperror.Internal(fmt.Errorf("panic: %v", rec)))
				c.AbortWithStatusJSON(status, env)
			}
		}()
		c.Next()
	}
}
