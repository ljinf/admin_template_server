package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/handler"
	model "poem_server_admin/internal/model/sys"
	service "poem_server_admin/internal/service/sys"
	"poem_server_admin/pkg/md5"
	"strconv"
	"time"
)

type UserHandler interface {
	CreateUser(ctx *gin.Context)
	UserList(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	Profile(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DelUser(ctx *gin.Context)
	AuthRoleById(ctx *gin.Context) //根据用户编号获取授权角色
	AuthRole(ctx *gin.Context)     //用户授权角色
	DeptTree(ctx *gin.Context)     //获取部门树列表
	//更新用户信息
	UpdateProfile(ctx *gin.Context)
}

type userHandler struct {
	*handler.Handler
	userService service.UserService
	roleService service.RoleService
}

func NewUserHandler(handler *handler.Handler, userService service.UserService, roleService service.RoleService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
		roleService: roleService,
	}
}

func (h *userHandler) UserList(ctx *gin.Context) {

	params := handler.GetParams(ctx, []string{"status", "username", "phone_number", "page_num", "page_size"})

	userList, err := h.userService.UserList(ctx, params["username"], params["phone_number"], params["status"],
		params["page_num"], params["page_size"])
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
	}

	v1.HandleSuccess(ctx, userList)
}

func (h *userHandler) CreateUser(ctx *gin.Context) {
	var param v1.EditUserRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	user := model.SysUser{
		DeptId:      param.DeptId,
		Username:    param.Username,
		NickName:    param.NickName,
		Email:       param.Email,
		PhoneNumber: param.PhoneNumber,
		Password:    param.Password,
		Sex:         param.Sex,
		Avatar:      param.Avatar,
		Status:      param.Status,
		LoginIp:     ctx.RemoteIP(),
		LoginDate:   time.Now(),
	}
	userId, err := h.userService.AddUser(ctx, &user)
	if err != nil {
		h.Logger.Error(err.Error())
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				v1.HandleError(ctx, http.StatusOK, v1.ErrUsernameAlreadyUse, nil)
				return
			}
		}
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	if len(param.RoleIds) > 0 {
		ru := make([]model.SysRoleUser, 0, len(param.RoleIds))
		for _, v := range param.RoleIds {
			ru = append(ru, model.SysRoleUser{
				UserId: userId,
				RoleId: int64(v),
			})
		}
		if err := h.roleService.UpdateUserRole(ctx, ru...); err != nil {
			h.Logger.Error(err.Error())
		}
	}

	v1.HandleSuccess(ctx, nil)

}

func (h *userHandler) AuthRoleById(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *userHandler) AuthRole(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *userHandler) DeptTree(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *userHandler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("userId")
	userId, _ := strconv.Atoi(id)

	res := make(map[string]interface{})

	//角色信息
	roleList, err := h.roleService.GetRoleAll(ctx)
	if err != nil {
		h.Logger.Error(err.Error())
	}
	res["roles"] = []interface{}{}
	if len(roleList) > 0 {
		if userId == 1 {
			res["roles"] = roleList
		} else {
			var list []model.SysRole
			for _, v := range roleList {
				if v.Id != 1 {
					list = append(list, v)
				}
			}
			res["roles"] = list
		}
	}

	if userId != 0 {
		user, err := h.userService.GetUserById(ctx, int64(userId))
		if err != nil {
			h.Logger.Error(err.Error(), zap.Any("user", user))
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}
		res["data"] = user

		roleIds, err := h.roleService.GetRolesByUserId(ctx, int64(userId))
		if err != nil {
			h.Logger.Error(err.Error())
		}
		res["roleIds"] = roleIds
	}

	v1.HandleSuccess(ctx, res)
}

func (h *userHandler) Profile(ctx *gin.Context) {
	id, ok := ctx.Get("userId")
	if !ok {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}
	userId := id.(int64)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	user, err := h.userService.GetUserById(ctx, userId)
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}
	roleGroup := []string{}

	if userId == 1 {
		roleGroup = append(roleGroup, "超级管理员")
	} else {
		//角色信息
		roleList, err := h.roleService.GetRolePermissionByUserId(ctx, userId)
		if err != nil {
			h.Logger.Error(err.Error())
		}
		for _, v := range roleList {
			roleGroup = append(roleGroup, v.RoleName)
		}
	}

	res := map[string]interface{}{
		"user_info": user,
		"role":      roleGroup,
	}

	v1.HandleSuccess(ctx, res)
}

func (h *userHandler) UpdateUser(ctx *gin.Context) {
	var param v1.EditUserRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	user := model.SysUser{
		Id:          param.Id,
		DeptId:      param.DeptId,
		Username:    param.Username,
		NickName:    param.NickName,
		Email:       param.Email,
		PhoneNumber: param.PhoneNumber,
		Password:    param.Password,
		Sex:         param.Sex,
		Avatar:      param.Avatar,
		Status:      param.Status,
		LoginIp:     ctx.RemoteIP(),
		LoginDate:   time.Now(),
	}

	if err := h.userService.UpdateUser(ctx, &user); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	if len(param.RoleIds) > 0 {
		ru := make([]model.SysRoleUser, 0, len(param.RoleIds))
		for _, v := range param.RoleIds {
			ru = append(ru, model.SysRoleUser{
				UserId: param.Id,
				RoleId: int64(v),
			})
		}
		if err := h.roleService.UpdateUserRole(ctx, ru...); err != nil {
			h.Logger.Error(err.Error())
		}
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *userHandler) DelUser(ctx *gin.Context) {
	id := ctx.Param("userId")
	userId, _ := strconv.Atoi(id)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if err := h.userService.DelUser(ctx, int64(userId)); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *userHandler) UpdateProfile(ctx *gin.Context) {

	id, ok := ctx.Get("userId")
	if !ok {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}
	userId := id.(int64)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	var param v1.UpdateProfileRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if param.OldPassword != "" {
		sysUser, err := h.userService.GetUserById(ctx, userId)
		if err != nil {
			h.Logger.Error(err.Error())
			v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
			return
		}

		if sysUser.Password != md5.Md5(param.OldPassword) {
			v1.HandleError(ctx, http.StatusOK, v1.ErrPwdError, nil)
			return
		}
	}

	user := model.SysUser{
		Id:          userId,
		NickName:    param.NickName,
		Email:       param.Email,
		PhoneNumber: param.PhoneNumber,
		Password:    param.Password,
		Sex:         param.Sex,
	}

	if err := h.userService.UpdateUser(ctx, &user); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
