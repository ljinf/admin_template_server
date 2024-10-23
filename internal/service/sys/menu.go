package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	model "poem_server_admin/internal/model/sys"
	repository "poem_server_admin/internal/repository/sys"
	"poem_server_admin/internal/service"
)

type MenuService interface {
	GetRouters(ctx context.Context, userId int64) ([]*model.SysMenu, error)
	MenuList(ctx context.Context, conds model.SysMenu) ([]*model.SysMenu, error)
	GetMenuInfoById(ctx context.Context, menuId int64) (*model.SysMenu, error)
	GetMenuInfoByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error)
	GetMenuInfoByRoleId(ctx context.Context, roleId int64) ([]*model.SysMenu, error)
	GetMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	AddMenu(ctx context.Context, menu *model.SysMenu) error
	UpdateMenu(ctx context.Context, menu *model.SysMenu) error
	DelMenu(ctx context.Context, menuId int64) error
}

type menuService struct {
	*service.Service
	menuRepository repository.MenuRepository
}

func NewMenuService(service *service.Service, repo repository.MenuRepository) MenuService {
	return &menuService{
		Service:        service,
		menuRepository: repo,
	}
}

func (m *menuService) GetMenuInfoByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error) {
	if m.isAdmin(userId) {
		return m.MenuList(ctx, model.SysMenu{Status: "0"})
	}
	return m.menuRepository.GetMenuByUserId(ctx, userId)
}

func (m *menuService) GetMenuInfoByRoleId(ctx context.Context, roleId int64) ([]*model.SysMenu, error) {
	return m.menuRepository.SelectMenuByRoleId(ctx, roleId)
}

func (m *menuService) GetMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	return m.menuRepository.GetMenuIdsByRoleId(ctx, roleId)
}

func (m *menuService) GetMenuInfoById(ctx context.Context, menuId int64) (*model.SysMenu, error) {
	return m.menuRepository.GetMenuById(ctx, menuId)
}

func (m *menuService) AddMenu(ctx context.Context, menu *model.SysMenu) error {

	//todo 权限校验

	err := m.menuRepository.InsertMenu(ctx, menu)

	return err

}

func (m *menuService) UpdateMenu(ctx context.Context, menu *model.SysMenu) error {
	//todo 权限校验

	err := m.menuRepository.EditMenu(ctx, menu)
	return err
}

func (m *menuService) DelMenu(ctx context.Context, menuId int64) error {
	//todo 权限校验
	//  warn("存在子菜单,不允许删除");
	// warn("菜单已分配,不允许删除");

	err := m.menuRepository.DeleteMenu(ctx, menuId)

	return err
}

func (m *menuService) MenuList(ctx context.Context, conds model.SysMenu) ([]*model.SysMenu, error) {
	list, err := m.menuRepository.MenuList(ctx, conds)
	if err != nil {
		m.Logger.Error(err.Error(), zap.Any("conds", conds))
		return nil, err
	}
	return list, nil
}

func (m *menuService) GetRouters(ctx context.Context, userId int64) ([]*model.SysMenu, error) {
	//权限校验，返回对应角色的菜单
	if userId == 1 {
		//管理员
		menus, err := m.menuRepository.SelectMenuTreeAll(ctx)
		if err != nil {
			m.Logger.Error(fmt.Sprintf("SelectMenuTreeAll %+v", err), zap.Any("userId", userId))
			return nil, err
		}

		return m.getChildPerms(menus, 0)
	}

	menus, err := m.menuRepository.SelectMenuTreeByUserId(ctx, userId)
	if err != nil {
		m.Logger.Error(fmt.Sprintf("SelectMenuTreeByUserId %+v", err), zap.Any("userId", userId))
		return nil, err
	}

	return m.getChildPerms(menus, 0)
}

// 根据父节点的ID获取所有子节点
func (m *menuService) getChildPerms(list []*model.SysMenu, parentId int64) ([]*model.SysMenu, error) {

	resList := []*model.SysMenu{}

	for _, v := range list {

		// 一、根据传入的某个父节点ID,遍历该父节点的所有子节点
		if v.ParentId == parentId {
			m.recursionFn(list, v)
			resList = append(resList, v)
		}
	}
	return resList, nil
}

// 递归列表
func (m *menuService) recursionFn(list []*model.SysMenu, t *model.SysMenu) {

	// 得到子节点列表
	childList := m.getChildList(list, t)
	t.Children = childList
	for _, tChild := range childList {
		if m.hasChild(list, tChild) {
			m.recursionFn(list, tChild)
		}
	}
}

// 得到子节点列表
func (m *menuService) getChildList(list []*model.SysMenu, t *model.SysMenu) []*model.SysMenu {

	tlist := []*model.SysMenu{}
	for _, v := range list {
		if v.ParentId == t.MenuId {
			tlist = append(tlist, v)
		}
	}
	return tlist

}

// 判断是否有子节点
func (m *menuService) hasChild(list []*model.SysMenu, t *model.SysMenu) bool {
	return len(m.getChildList(list, t)) > 0
}

// 超级管理员
func (m *menuService) isAdmin(userId int64) bool {
	return userId == 1
}
