package sys

import (
	"context"
	model "poem_server_admin/internal/model/sys"
	"poem_server_admin/internal/repository"
	"strings"
)

type UserRepository interface {
	FirstById(ctx context.Context, id int64) (*model.SysUser, error)
	UserInfoByName(ctx context.Context, name string) (*model.SysUser, error)
	SelectAllUsers(ctx context.Context, username, phone, status, pageNum, pageSize string) (map[string]interface{}, error)
	CreateUser(ctx context.Context, user *model.SysUser) (int64, error)
	EditUser(ctx context.Context, user *model.SysUser) error
	DeleteUser(ctx context.Context, userid int64) error
}
type userRepository struct {
	*repository.Repository
}

func NewUserRepository(repository *repository.Repository) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

func (r *userRepository) FirstById(ctx context.Context, id int64) (*model.SysUser, error) {
	user := model.SysUser{Id: id}
	if err := r.DB(ctx).Where("id=?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UserInfoByName(ctx context.Context, name string) (*model.SysUser, error) {
	user := model.SysUser{}
	if err := r.DB(ctx).Where("username=?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SelectAllUsers(ctx context.Context, username, phone, status, pageNum, pageSize string) (map[string]interface{}, error) {
	var (
		conds  = []string{}
		values = []interface{}{}
	)

	if username != "" {
		conds = append(conds, "username like ?")
		values = append(values, "%"+username+"%")
	}

	if phone != "" {
		conds = append(conds, "phone_number = ?")
		values = append(values, phone)
	}

	if status != "" {
		conds = append(conds, "status = ?")
		values = append(values, status)
	}
	page, size := repository.ParsePage(pageNum, pageSize)

	var list []model.SysUser
	if err := r.DB(ctx).Select("id, username,nick_name,user_type,email,phone_number,sex,avatar,status,created_at").
		Where(strings.Join(conds, " and "), values...).Limit(size).
		Offset((page - 1) * size).Find(&list).Error; err != nil {
		return nil, err
	}
	var count int64
	r.DB(ctx).Table("sys_user").Where(strings.Join(conds, " and "), values...).Count(&count)

	return map[string]interface{}{
		"total": count,
		"rows":  list,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.SysUser) (int64, error) {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *userRepository) EditUser(ctx context.Context, user *model.SysUser) error {
	return r.DB(ctx).Where("id=?", user.Id).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, userid int64) error {
	return r.DB(ctx).Where("id=?", userid).Delete(&model.SysUser{}).Error
}
