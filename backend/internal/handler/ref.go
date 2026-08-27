package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

type RefHandler struct {
	svc *service.RefService
}

func NewRefHandler(svc *service.RefService) *RefHandler {
	return &RefHandler{svc: svc}
}

func (h *RefHandler) Provinces(c *gin.Context) {
	items, err := h.svc.Provinces(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	writeCached(c, items)
}

func (h *RefHandler) Districts(c *gin.Context) {
	code := c.Query("province_code")
	items, err := h.svc.Districts(c.Request.Context(), code)
	if err != nil {
		writeErr(c, err)
		return
	}
	writeCached(c, items)
}

func (h *RefHandler) Subdistricts(c *gin.Context) {
	code := c.Query("district_code")
	items, err := h.svc.Subdistricts(c.Request.Context(), code)
	if err != nil {
		writeErr(c, err)
		return
	}
	writeCached(c, items)
}

func writeCached(c *gin.Context, data any) {
	env := httpx.OK(data)
	body, err := json.Marshal(env)
	if err != nil {
		writeErr(c, err)
		return
	}
	sum := sha256.Sum256(body)
	etag := `"` + hex.EncodeToString(sum[:8]) + `"`
	c.Header("ETag", etag)
	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", body)
}
