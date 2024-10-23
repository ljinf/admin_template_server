package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/handler"
	service "poem_server_admin/internal/service/sys"
)

type AccountHandler interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetInfo(ctx *gin.Context)
}

type accountHandler struct {
	*handler.Handler
	userService service.UserService
	roleService service.RoleService
	menuService service.MenuService
}

func NewAccountHandler(handler *handler.Handler, us service.UserService, rs service.RoleService, ms service.MenuService) AccountHandler {
	return &accountHandler{
		Handler:     handler,
		userService: us,
		roleService: rs,
		menuService: ms,
	}
}

func (h *accountHandler) Login(ctx *gin.Context) {

	var req v1.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Logger.Error(err.Error(), zap.String("loginInfo", fmt.Sprintf("%v", req)))
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	token, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrLoginError, nil)
		return
	}

	//创建token
	res := v1.LoginResp{
		AccessToken: token,
	}

	v1.HandleSuccess(ctx, res)
}

func (h *accountHandler) GetInfo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	user, err := h.userService.GetUserById(ctx, userId.(int64))
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	user.Password = ""

	var resRole []string  //角色
	var resPerms []string //权限字符

	var roleIds []int64

	if h.isAdmin(userId.(int64)) {
		resRole = []string{"admin"}
		resPerms = []string{"*:*:*"}
	} else {

		roles, err := h.roleService.GetRolePermissionByUserId(ctx, userId.(int64))
		if err != nil {
			v1.HandleError(ctx, http.StatusOK, v1.ErrRolePermsFailed, nil)
			return
		}

		resRole = make([]string, 0, len(roles))
		roleIds = make([]int64, 0, len(roles))
		for _, v := range roles {
			roleIds = append(roleIds, v.Id)
			resRole = append(resRole, v.RoleKey)
		}

		if len(roleIds) > 0 {
			for _, v := range roleIds {
				if menus, err := h.menuService.GetMenuInfoByRoleId(ctx, v); err != nil {
					h.Logger.Error(err.Error())
				} else {
					for _, m := range menus {
						resPerms = append(resPerms, m.Perms)
					}
				}
			}
		}
	}

	res := map[string]interface{}{
		"user":        user,
		"roles":       resRole,
		"permissions": resPerms,
	}

	v1.HandleSuccess(ctx, res)
}

func (h *accountHandler) Logout(ctx *gin.Context) {
	v1.HandleSuccess(ctx, nil)
}

func (h *accountHandler) isAdmin(userId int64) bool {
	return userId == 1
}
