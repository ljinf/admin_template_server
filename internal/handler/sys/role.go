package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/handler"
	model "poem_server_admin/internal/model/sys"
	service "poem_server_admin/internal/service/sys"
	"strconv"
)

type RoleHandler interface {
	RoleList(ctx *gin.Context)
	RoleInfo(ctx *gin.Context)
	AddRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DelRole(ctx *gin.Context)
}

type roleHandler struct {
	*handler.Handler
	roleService service.RoleService
}

func NewRoleHandler(h *handler.Handler, srv service.RoleService) RoleHandler {
	return &roleHandler{
		Handler:     h,
		roleService: srv,
	}
}

func (h *roleHandler) RoleList(ctx *gin.Context) {
	roleList, err := h.roleService.GetRoleList(ctx)
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, roleList)
}

func (h *roleHandler) RoleInfo(ctx *gin.Context) {
	id := ctx.Param("roleId")
	roleId, _ := strconv.Atoi(id)
	if roleId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}
	sysRole, err := h.roleService.GetRoleById(ctx, int64(roleId))
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, sysRole)

}

func (h *roleHandler) AddRole(ctx *gin.Context) {
	var param v1.RoleRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	req := model.SysRole{
		RoleName:  param.RoleName,
		RoleKey:   param.RoleKey,
		RoleSort:  param.RoleSort,
		DataScope: param.DataScope,
		Status:    param.Status,
	}
	roleId, err := h.roleService.AddRole(ctx, req)
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	var rm []model.SysRoleMenu
	if len(param.MenuIds) > 0 {
		for _, v := range param.MenuIds {
			item := model.SysRoleMenu{
				RoleId: roleId,
				MenuId: int64(v),
			}
			rm = append(rm, item)
		}

		if err := h.roleService.AddRoleMenu(ctx, rm...); err != nil {
			h.Logger.Error(err.Error())
			v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
			return
		}
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *roleHandler) UpdateRole(ctx *gin.Context) {
	var param v1.RoleRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if param.RoleId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	req := model.SysRole{
		Id:        param.RoleId,
		RoleName:  param.RoleName,
		RoleKey:   param.RoleKey,
		RoleSort:  param.RoleSort,
		DataScope: param.DataScope,
		Status:    param.Status,
	}

	if err := h.roleService.UpdateRole(ctx, req); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	var rm []model.SysRoleMenu
	if len(param.MenuIds) > 0 {
		for _, v := range param.MenuIds {
			item := model.SysRoleMenu{
				RoleId: req.Id,
				MenuId: int64(v),
			}
			rm = append(rm, item)
		}

		if err := h.roleService.UpdateRoleMenu(ctx, rm...); err != nil {
			h.Logger.Error(err.Error())
			v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
			return
		}
	}

	v1.HandleSuccess(ctx, nil)

}

func (h *roleHandler) DelRole(ctx *gin.Context) {
	id := ctx.Param("roleId")
	roleId, _ := strconv.Atoi(id)
	if roleId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if err := h.roleService.DelRole(ctx, int64(roleId)); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
