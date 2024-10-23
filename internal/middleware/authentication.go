package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/pkg/jwt"
)

func AuthMiddleware(jwt *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			v1.HandleError(ctx, http.StatusNonAuthoritativeInfo, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
