package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/cache"
	"poem_server_admin/pkg/log"
)

func CheckPerm(perm string, logger *log.Logger, accountCache cache.AccountCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := ctx.Get("userId")
		if !ok {
			v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		if userId.(int64) == 1 {
			ctx.Next()
			return
		}

		permissions, err := accountCache.GetUserPermissions(userId.(int64))
		if err != nil {
			logger.Error(err.Error())
		}

		if len(permissions) > 0 {
			for _, v := range permissions {
				if v == perm {
					ctx.Next()
					return
				}
			}
		}

		v1.HandleError(ctx, http.StatusOK, v1.ErrForbid, nil)
		ctx.Abort()
	}
}
