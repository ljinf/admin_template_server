package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"poem_server_admin/pkg/jwt"
	"poem_server_admin/pkg/log"
	"strconv"
	"time"
)

type Handler struct {
	Logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		Logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) int64 {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	return v.(*jwt.MyCustomClaims).UserId
}

func GetParams(ctx *gin.Context, keys []string) map[string]string {
	res := make(map[string]string)
	if len(keys) > 0 {
		for _, v := range keys {
			//begin_time end_time
			if v == "begin_time" || v == "end_time" {
				if parse, err := time.Parse("2006-01-02", ctx.Query(v)); err == nil {
					var day = parse.Day()
					if v == "end_time" {
						day = parse.Day() + 1
					}
					res[v] = fmt.Sprintf("%v", time.Date(parse.Year(), parse.Month(), day, 0, 0, 0, 0, time.Local).Unix())
				} else if _, err = strconv.Atoi(ctx.Query(v)); err == nil {
					res[v] = ctx.Query(v)
				}
			} else {
				res[v] = ctx.Query(v)
			}
		}
	}
	return res
}
