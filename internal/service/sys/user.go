package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/cache"
	model "poem_server_admin/internal/model/sys"
	repository "poem_server_admin/internal/repository/sys"
	"poem_server_admin/internal/service"
	"poem_server_admin/pkg/md5"
	"time"
)

const (
	Expired = 86400
)

type UserService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetUserById(ctx context.Context, id int64) (*model.SysUser, error)
	GetUserByUserName(ctx context.Context, name string) (*model.SysUser, error)

	UserList(ctx context.Context, username, phone, status, pageNum, pageSize string) (interface{}, error)
	AddUser(ctx context.Context, user *model.SysUser) (int64, error)
	UpdateUser(ctx context.Context, user *model.SysUser) error
	DelUser(ctx context.Context, userId int64) error
}

type userService struct {
	*service.Service
	userRepository repository.UserRepository
	userCache      cache.AccountCache
}

func NewUserService(service *service.Service, userRepository repository.UserRepository,
	userCache cache.AccountCache) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
		userCache:      userCache,
	}
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	sysUser, err := s.GetUserByUserName(ctx, username)
	if err != nil {
		s.Logger.Error(err.Error(), zap.String("GetUserByUserName", fmt.Sprintf("%v", username)))
		return "", err
	}

	//校验密码
	if sysUser.Password != md5.Md5(password) {
		s.Logger.Error("密码错误", zap.String("login err", fmt.Sprintf("name:%v", username)))
		return "", errors.New("账号密码错误")
	}

	token, err := s.Jwt.GenToken(sysUser.Id, time.Unix(time.Now().Unix()+int64(Expired), 0))
	if err != nil {
		s.Logger.Error(err.Error())
		return "", err
	}

	if err = s.userCache.SetTokenCache(sysUser.Id, token, Expired); err != nil {
		s.Logger.Error(err.Error(), zap.String("userId", fmt.Sprintf("%v", sysUser.Id)))
	}

	//if err = s.userCache.SetUserInfoCache(sysUser, Expired); err != nil {
	//	s.Logger.Error(err.Error(), zap.Any("user", sysUser))
	//}

	return token, nil

}

func (s *userService) GetUserByUserName(ctx context.Context, name string) (*model.SysUser, error) {
	sysUser, err := s.userRepository.UserInfoByName(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error(v1.ErrUserNotFoundError.Error(), zap.String("username", name))
		} else {
			s.Logger.Error(err.Error(), zap.String("username", name))
		}
	}
	return sysUser, err
}

func (s *userService) GetUserById(ctx context.Context, id int64) (*model.SysUser, error) {

	userInfoCache, err := s.userCache.GetUserInfoCache(id)
	if err != nil {
		s.Logger.Error(err.Error(), zap.Any("userId", id))
	}

	if userInfoCache != nil {
		return userInfoCache["user"].(*model.SysUser), nil
	}

	sysUser, err := s.userRepository.FirstById(ctx, id)
	if err != nil {
		s.Logger.Error(err.Error(), zap.Any("userId", id))
		return nil, err
	}

	return sysUser, nil
}

func (s *userService) UserList(ctx context.Context, username, phone, status, pageNum, pageSize string) (interface{}, error) {
	users, err := s.userRepository.SelectAllUsers(ctx, username, phone, status, pageNum, pageSize)
	return users, err
}

func (s *userService) AddUser(ctx context.Context, user *model.SysUser) (int64, error) {
	user.Password = md5.Md5(user.Password)
	return s.userRepository.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.SysUser) error {
	if user.Password != "" {
		user.Password = md5.Md5(user.Password)
	}
	return s.userRepository.EditUser(ctx, user)
}

func (s *userService) DelUser(ctx context.Context, userId int64) error {
	return s.userRepository.DeleteUser(ctx, userId)
}
