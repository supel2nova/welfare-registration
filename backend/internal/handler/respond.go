package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

func writeOK(c *gin.Context, status int, data any) {
	c.JSON(status, httpx.OK(data))
}

func writeErr(c *gin.Context, err error) {
	var ae *apperror.Error
	if !errors.As(err, &ae) {
		ae = apperror.Internal(err)
	}
	if ae.HTTPStatus >= 500 {
		log.Printf("%s %s -> %d %v", c.Request.Method, c.Request.URL.Path, ae.HTTPStatus, ae)
	}
	c.JSON(httpx.Fail(ae))
}
