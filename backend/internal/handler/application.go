package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/middleware"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
)

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeErr(c, apperror.BadRequest(apperror.CodeInvalidPayload))
		return
	}

	res, err := h.svc.Create(c.Request.Context(), middleware.ActorFrom(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	writeOK(c, http.StatusCreated, res)
}
