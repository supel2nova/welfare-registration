package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
)

type DevHandler struct {
	svc *service.RefService
}

func NewDevHandler(svc *service.RefService) *DevHandler {
	return &DevHandler{svc: svc}
}

func (h *DevHandler) Users(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	writeOK(c, http.StatusOK, users)
}
