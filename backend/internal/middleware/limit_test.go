package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/supel2nova/welfare-registration/backend/internal/middleware"
)

func TestInFlightShedsWhenFull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	release := make(chan struct{})
	entered := make(chan struct{}, 1)

	r := gin.New()
	r.Use(middleware.InFlight(1))
	r.GET("/x", func(c *gin.Context) {
		entered <- struct{}{}
		<-release
		c.String(http.StatusOK, "ok")
	})

	first := httptest.NewRecorder()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.ServeHTTP(first, httptest.NewRequest("GET", "/x", nil))
	}()

	<-entered // สล็อตเดียวถูกจองแล้ว

	second := httptest.NewRecorder()
	r.ServeHTTP(second, httptest.NewRequest("GET", "/x", nil))
	if second.Code != http.StatusServiceUnavailable {
		t.Fatalf("คำขอที่สองควรได้ 503 แต่ได้ %d", second.Code)
	}
	if !strings.Contains(second.Body.String(), "SYS002") {
		t.Fatalf("ควรมี errorCode SYS002: %s", second.Body.String())
	}
	if second.Header().Get("Retry-After") == "" {
		t.Fatal("ควรมี header Retry-After")
	}

	close(release)
	wg.Wait()
	if first.Code != http.StatusOK {
		t.Fatalf("คำขอแรกควรได้ 200 แต่ได้ %d", first.Code)
	}

	third := httptest.NewRecorder()
	r.ServeHTTP(third, httptest.NewRequest("GET", "/x", nil))
	if third.Code != http.StatusOK {
		t.Fatalf("สล็อตต้องถูกคืนหลังคำขอแรกจบ แต่ได้ %d", third.Code)
	}
}

func TestBodyLimitRejectsOversizedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.BodyLimit(64))
	r.POST("/x", func(c *gin.Context) {
		var body map[string]any
		if err := c.ShouldBindJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "rejected")
			return
		}
		c.String(http.StatusOK, "ok")
	})

	small := httptest.NewRecorder()
	r.ServeHTTP(small, httptest.NewRequest("POST", "/x", strings.NewReader(`{"a":1}`)))
	if small.Code != http.StatusOK {
		t.Fatalf("payload เล็กต้องผ่าน แต่ได้ %d", small.Code)
	}

	big := httptest.NewRecorder()
	r.ServeHTTP(big, httptest.NewRequest("POST", "/x", strings.NewReader(`{"a":"`+strings.Repeat("x", 500)+`"}`)))
	if big.Code != http.StatusBadRequest {
		t.Fatalf("payload เกินลิมิตต้องถูกปฏิเสธ แต่ได้ %d", big.Code)
	}
}
