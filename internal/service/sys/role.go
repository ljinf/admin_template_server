package service

import (
	"context"
	model "poem_server_admin/internal/model/sys"
	repository "poem_server_admin/internal/repository/sys"
	"poem_server_admin/internal/service"
)

type RoleService interface {
	GetRoleList(ctx context.Context) (interface{}, error)
	GetRoleAll(ctx context.Context) ([]model.SysRole, error)
	GetRoleById(ctx context.Context, roleId int64) (*model.SysRole, error)
	AddRole(ctx context.Context, role model.SysRole) (int64, error)
	AddRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error
	UpdateRole(ctx context.Context, role model.SysRole) error
	UpdateRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error
	DelRole(ctx context.Context, roleId int64) error

	//用户角色
	GetRolesByUserId(ctx context.Context, userId int64) ([]int64, error)
	UpdateUserRole(ctx context.Context, roelUser ...model.SysRoleUser) error
	GetRolePermissionByUserId(ctx context.Context, userId int64) ([]model.SysRole, error)
}

type roleService struct {
	*service.Service
	repo repository.RoleRepository
}

func NewRoleService(srv *service.Service, repo repository.RoleRepository) RoleService {
	return &roleService{
		Service: srv,
		repo:    repo,
	}
}

func (r *roleService) GetRoleList(ctx context.Context) (interface{}, error) {
	//todo 分页
	roleList, err := r.repo.RoleList(ctx)
	return roleList, err
}

func (r *roleService) GetRoleAll(ctx context.Context) ([]model.SysRole, error) {
	return r.repo.RoleAll(ctx)
}

func (r *roleService) GetRoleById(ctx context.Context, roleId int64) (*model.SysRole, error) {
	roleInfo, err := r.repo.RoleInfo(ctx, roleId)
	return roleInfo, err
}

func (r *roleService) AddRole(ctx context.Context, role model.SysRole) (int64, error) {
	err := r.repo.InsertRole(ctx, &role)
	return role.Id, err
}

func (r *roleService) AddRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error {
	return r.repo.InsertRoleMenu(ctx, roleMenu...)
}

func (r *roleService) UpdateRole(ctx context.Context, role model.SysRole) error {
	err := r.repo.EditRole(ctx, &role)
	return err
}

func (r *roleService) UpdateRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error {
	err := r.repo.EditRoleMenu(ctx, roleMenu...)
	return err
}

func (r *roleService) DelRole(ctx context.Context, roleId int64) error {
	return r.Tm.Transaction(ctx, func(ctx context.Context) error {
		if err := r.repo.DeleteRoleMenu(ctx, roleId); err != nil {
			return err
		}
		return r.repo.DeleteRole(ctx, roleId)
	})
}

func (r *roleService) GetRolesByUserId(ctx context.Context, userId int64) ([]int64, error) {
	return r.repo.SelectRolesByUserId(ctx, userId)
}

// 更新用户角色
func (r *roleService) UpdateUserRole(ctx context.Context, roelUser ...model.SysRoleUser) error {
	return r.repo.EditRoleUser(ctx, roelUser...)
}

// 获取角色的权限
func (r *roleService) GetRolePermissionByUserId(ctx context.Context, userId int64) ([]model.SysRole, error) {

	roles, err := r.repo.SelectRolePermissionByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleService) isAdmin(userId int64) bool {
	return userId == 1
}
