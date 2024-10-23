package sys

import (
	"context"
	"gorm.io/gorm"
	model "poem_server_admin/internal/model/sys"
	"poem_server_admin/internal/repository"
)

type RoleRepository interface {
	RoleList(ctx context.Context) (map[string]interface{}, error)
	RoleAll(ctx context.Context) ([]model.SysRole, error)
	RoleInfo(ctx context.Context, roleId int64) (*model.SysRole, error)
	InsertRole(ctx context.Context, role *model.SysRole) error
	InsertRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error //角色菜单
	EditRoleUser(ctx context.Context, roleUser ...model.SysRoleUser) error   //角色用户
	EditRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error
	EditRole(ctx context.Context, role *model.SysRole) error
	DeleteRole(ctx context.Context, roleId int64) error
	DeleteRoleMenu(ctx context.Context, roleId int64) error

	//用户角色
	SelectRolesByUserId(ctx context.Context, userId int64) ([]int64, error)
	SelectRolePermissionByUserId(ctx context.Context, userId int64) ([]model.SysRole, error)
}

type roleRepository struct {
	*repository.Repository
}

func NewRoleRepository(repo *repository.Repository) RoleRepository {
	return &roleRepository{
		Repository: repo,
	}
}

func (r *roleRepository) RoleList(ctx context.Context) (map[string]interface{}, error) {
	var list []model.SysRole
	err := r.DB(ctx).Find(&list).Error
	var count int64
	r.DB(ctx).Table("sys_role").Count(&count)
	return map[string]interface{}{
		"total": count,
		"rows":  list,
	}, err
}

func (r *roleRepository) RoleAll(ctx context.Context) ([]model.SysRole, error) {
	var list []model.SysRole
	err := r.DB(ctx).Find(&list).Error
	return list, err
}

func (r *roleRepository) RoleInfo(ctx context.Context, roleId int64) (*model.SysRole, error) {
	role := model.SysRole{Id: roleId}
	err := r.DB(ctx).First(&role).Error
	return &role, err
}

func (r *roleRepository) InsertRole(ctx context.Context, role *model.SysRole) error {
	err := r.DB(ctx).Create(role).Error
	return err
}

func (r *roleRepository) EditRole(ctx context.Context, role *model.SysRole) error {
	err := r.DB(ctx).Save(role).Error
	return err
}

func (r *roleRepository) DeleteRole(ctx context.Context, roleId int64) error {
	return r.DB(ctx).Delete(&model.SysRole{
		Id: roleId,
	}).Error
}

func (r *roleRepository) InsertRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error {
	err := r.DB(ctx).Create(&roleMenu).Error
	return err
}

func (r *roleRepository) EditRoleMenu(ctx context.Context, roleMenu ...model.SysRoleMenu) error {
	if len(roleMenu) > 0 {
		roleId := roleMenu[0].RoleId

		r.DB(ctx).Transaction(func(tx *gorm.DB) error {

			if err := tx.Where("role_id=?", roleId).Delete(&model.SysRoleMenu{}).Error; err != nil {
				return err
			}

			return tx.Create(&roleMenu).Error
		})
	}
	return nil
}

func (r *roleRepository) DeleteRoleMenu(ctx context.Context, roleId int64) error {
	return r.DB(ctx).Where("role_id=?", roleId).Delete(&model.SysRoleMenu{}).Error
}

func (r *roleRepository) SelectRolesByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var list []int64
	err := r.DB(ctx).Table("sys_user_role").Where("user_id=?", userId).Select("role_id").Find(&list).Error
	return list, err
}

// 更新用户角色，先删除，在插入
func (r *roleRepository) EditRoleUser(ctx context.Context, roleUser ...model.SysRoleUser) error {

	if len(roleUser) > 0 {
		return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
			if err := r.DB(ctx).Table("sys_user_role").Where("user_id=?", roleUser[0].UserId).Delete(&model.SysRoleUser{}).Error; err != nil {
				return err
			}
			return r.DB(ctx).Create(&roleUser).Error
		})
	}

	return nil
}

func (r *roleRepository) SelectRolePermissionByUserId(ctx context.Context, userId int64) ([]model.SysRole, error) {

	querySql := "select distinct r.id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly," +
		"r.status, r.del_flag, r.create_time, r.remark " +
		"from sys_role r " +
		"left join sys_user_role ur on ur.role_id = r.id " +
		"left join sys_user u on u.id = ur.user_id " +
		"where r.del_flag = '0' and ur.user_id = ?"

	list := []model.SysRole{}
	err := r.DB(ctx).Raw(querySql, userId).Scan(&list).Error
	return list, err
}
