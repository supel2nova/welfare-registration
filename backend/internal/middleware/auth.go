package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

const (
	headerDebugUser = "X-Debug-User-Id"
	actorKey        = "actor"
)

func StubAuth(repo *repository.Repo, enabled bool, defaultUserID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			status, env := httpx.Fail(apperror.Unauthorized())
			c.AbortWithStatusJSON(status, env)
			return
		}

		raw := c.GetHeader(headerDebugUser)
		if raw == "" {
			raw = defaultUserID
		}
		if raw == "" {
			status, env := httpx.Fail(apperror.Unauthorized())
			c.AbortWithStatusJSON(status, env)
			return
		}

		id, err := uuid.Parse(raw)
		if err != nil {
			status, env := httpx.Fail(apperror.Unauthorized())
			c.AbortWithStatusJSON(status, env)
			return
		}

		user, err := repo.FindUserByID(c.Request.Context(), id)
		if errors.Is(err, repository.ErrUserNotFound) {
			status, env := httpx.Fail(apperror.Unauthorized())
			c.AbortWithStatusJSON(status, env)
			return
		}
		if err != nil {
			status, env := httpx.Fail(apperror.Internal(err))
			c.AbortWithStatusJSON(status, env)
			return
		}

		uid := user.ID
		c.Set(actorKey, domain.Actor{
			Type:   domain.ActorUser,
			OrgID:  user.OrgID,
			UserID: &uid,
			Role:   domain.Role(user.Role),
		})
		c.Next()
	}
}

func ActorFrom(c *gin.Context) domain.Actor {
	v, ok := c.Get(actorKey)
	if !ok {
		return domain.Actor{}
	}
	actor, _ := v.(domain.Actor)
	return actor
}
