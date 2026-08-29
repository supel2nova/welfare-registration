package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/supel2nova/welfare-registration/backend/internal/config"
	"github.com/supel2nova/welfare-registration/backend/internal/middleware"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/service"
)

type Deps struct {
	Cfg  config.Config
	Pool *pgxpool.Pool
	Repo *repository.Repo
	Apps *service.ApplicationService
	Ref  *service.RefService
}

func NewRouter(d Deps) *gin.Engine {
	if d.Cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(middleware.Recovery(), middleware.CORS(d.Cfg.CORSOrigin), gin.Logger())

	// health ต้องลงทะเบียนก่อน InFlight ไม่งั้นตอนโหลดพีค probe จะโดนปฏิเสธ แล้ว k8s รีสตาร์ต pod ซ้ำเติม
	health := NewHealthHandler(d.Pool)
	r.GET("/health", health.Get)

	r.Use(middleware.InFlight(d.Cfg.MaxInFlight), middleware.BodyLimit(d.Cfg.MaxBodyBytes))

	ref := NewRefHandler(d.Ref)
	apps := NewApplicationHandler(d.Apps)
	dev := NewDevHandler(d.Ref)

	v1 := r.Group("/api/v1")
	{
		refGroup := v1.Group("/ref")
		refData := middleware.CacheControl(middleware.CacheRefData)
		refGroup.GET("/provinces", refData, ref.Provinces)
		refGroup.GET("/districts", refData, ref.Districts)
		refGroup.GET("/subdistricts", refData, ref.Subdistricts)
		refGroup.GET("/address-search", middleware.CacheControl(middleware.CacheSearch), ref.SearchAddress)

		auth := middleware.StubAuth(d.Repo, d.Cfg.StubAuthEnabled && !d.Cfg.IsProduction(), d.Cfg.StubDefaultUserID)
		v1.POST("/applications", auth, apps.Create)

		if d.Cfg.StubAuthEnabled && !d.Cfg.IsProduction() {
			v1.GET("/dev/users", dev.Users)
		}
	}

	return r
}
