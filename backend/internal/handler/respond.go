package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

func writeOK(c *gin.Context, status int, data any) {
	c.JSON(status, httpx.OK(data))
}

func writeErr(c *gin.Context, err error) {
	var ae *apperror.Error
	if errors.As(err, &ae) {
		c.JSON(httpx.Fail(ae))
		return
	}
	c.JSON(httpx.Fail(apperror.Internal(err)))
}
