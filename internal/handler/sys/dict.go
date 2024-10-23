package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/handler"
	service "poem_server_admin/internal/service/sys"
)

type DictHandler interface {
	GetDictDataType(ctx *gin.Context)
}

type dictHandler struct {
	*handler.Handler
	dictService service.DictService
}

func NewDictHandler(h *handler.Handler, dictService service.DictService) DictHandler {
	return &dictHandler{
		Handler:     h,
		dictService: dictService,
	}
}

func (h *dictHandler) GetDictDataType(ctx *gin.Context) {
	dictType := ctx.Param("dictType")
	if dictType == "" {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	dataTypeList, err := h.dictService.GetDictDataType(ctx, dictType)
	if err != nil {
		h.Logger.Error(err.Error(), zap.Any("dictType", dictType))
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, dataTypeList)
}
