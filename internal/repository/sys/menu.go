package sys

import (
	"context"
	model "poem_server_admin/internal/model/sys"
	"poem_server_admin/internal/repository"
	"strings"
)

type MenuRepository interface {
	InsertMenu(ctx context.Context, menu *model.SysMenu) error
	EditMenu(ctx context.Context, menu *model.SysMenu) error
	DeleteMenu(ctx context.Context, menuId int64) error
	GetMenuById(ctx context.Context, menuId int64) (*model.SysMenu, error)
	GetMenuByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error)
	GetMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	GetMenuByNameAndParentId(ctx context.Context, name string, parentId int64) (*model.SysMenu, error)
	MenuList(ctx context.Context, conds model.SysMenu) ([]*model.SysMenu, error)
	SelectMenuTreeAll(ctx context.Context) ([]*model.SysMenu, error)
	SelectMenuTreeByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error)
	SelectMenuByRoleId(ctx context.Context, roleId int64) ([]*model.SysMenu, error)
}

type menuRepository struct {
	*repository.Repository
}

func NewMenuRepository(repo *repository.Repository) MenuRepository {
	return &menuRepository{
		Repository: repo,
	}
}

func (m *menuRepository) GetMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	querySql := "select m.menu_id " +
		"from sys_menu m " +
		"left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"where rm.role_id = ? order by m.parent_id, m.order_num"
	var list []int64
	err := m.DB(ctx).Raw(querySql, roleId).Scan(&list).Error
	return list, err
}

func (m *menuRepository) GetMenuByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error) {
	querySql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.created_at " +
		"from sys_menu m " +
		"left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"left join sys_user_role ur on rm.role_id = ur.role_id " +
		"left join sys_role ro on ur.role_id = ro.id " +
		"where ur.user_id = ? and m.visible = 0 and m.status = 0 order by m.parent_id, m.order_num"

	var list []*model.SysMenu
	err := m.DB(ctx).Raw(querySql, userId).Scan(&list).Error
	return list, err
}

func (m *menuRepository) GetMenuById(ctx context.Context, menuId int64) (*model.SysMenu, error) {
	querySql := " select menu_id, menu_name, parent_id, order_num, path, component, `query`, is_frame, is_cache, menu_type, visible, status, ifnull(perms,'') as perms, icon, created_at " +
		"from sys_menu where menu_id = ?"

	var info model.SysMenu
	err := m.DB(ctx).Raw(querySql, menuId).Scan(&info).Error
	return &info, err
}

func (m *menuRepository) GetMenuByNameAndParentId(ctx context.Context, name string, parentId int64) (*model.SysMenu, error) {
	querySql := " select menu_id, menu_name, parent_id, order_num, path, component, `query`, is_frame, is_cache, menu_type, visible, status, ifnull(perms,'') as perms, icon, created_at " +
		"from sys_menu where menu_name=? and parent_id = ? limit 1"

	var info model.SysMenu
	err := m.DB(ctx).Raw(querySql, name, parentId).Scan(&info).Error
	return &info, err
}

func (m *menuRepository) InsertMenu(ctx context.Context, menu *model.SysMenu) error {
	return m.DB(ctx).Create(menu).Error
}

func (m *menuRepository) EditMenu(ctx context.Context, menu *model.SysMenu) error {
	return m.DB(ctx).Where("menu_id=?", menu.MenuId).Save(menu).Error
}

func (m *menuRepository) DeleteMenu(ctx context.Context, menuId int64) error {
	deleteSql := "delete from `sys_menu` where menu_id=?"
	return m.DB(ctx).Exec(deleteSql, menuId).Error
}

func (m *menuRepository) MenuList(ctx context.Context, conds model.SysMenu) ([]*model.SysMenu, error) {

	var (
		where = []string{" status = ? "}
		args  = []any{conds.Status}
	)

	if conds.MenuName != "" {
		where = append(where, "menu_name like concat('%"+conds.MenuName+"%')")
		args = append(args, conds.MenuName)
	}
	if conds.Visible != "" {
		where = append(where, "visible = ?")
		args = append(args, conds.Visible)
	}

	querySql := "select menu_id, menu_name, parent_id, order_num, path, component, `query`, is_frame, is_cache, menu_type, visible, status, ifnull(perms,'') as perms, icon, created_at " +
		"from sys_menu where " + strings.Join(where, " and ") + " order by parent_id, order_num"
	return m.query(ctx, querySql, args...)
}

func (m *menuRepository) SelectMenuTreeAll(ctx context.Context) ([]*model.SysMenu, error) {
	querySql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.created_at " +
		"from sys_menu m " +
		"where m.menu_type in ('M', 'C') and m.status = 0 " +
		"order by m.parent_id, m.order_num"
	return m.query(ctx, querySql)
}

func (m *menuRepository) SelectMenuTreeByUserId(ctx context.Context, userId int64) ([]*model.SysMenu, error) {
	//querySql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.created_at " +
	//	"from sys_menu m " +
	//	"left join sys_role_menu rm on m.menu_id = rm.menu_id " +
	//	"left join sys_user_role ur on rm.role_id = ur.role_id " +
	//	"left join sys_role ro on ur.role_id = ro.id " +
	//	"left join sys_user u on ur.user_id = u.id " +
	//	"where u.id = ? and m.menu_type in ('M', 'C') and m.status = 0  AND ro.status = 0 " +
	//	"order by m.parent_id, m.order_num"

	querySql := "SELECT DISTINCT m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, IFNULL(m.perms,'') AS perms, m.is_frame," +
		"m.is_cache, m.menu_type, m.icon, m.order_num, m.created_at " +
		"FROM sys_menu m WHERE m.menu_type IN ('M', 'C') AND m.status = 0 AND " +
		"m.`menu_id` IN (SELECT rm.`menu_id` FROM `sys_role_menu` rm WHERE rm.`role_id` IN (SELECT ur.`role_id` FROM `sys_user_role` ur WHERE ur.`user_id`=?))"
	return m.query(ctx, querySql, userId)
}

func (m *menuRepository) SelectMenuByRoleId(ctx context.Context, roleId int64) ([]*model.SysMenu, error) {
	querySql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.created_at " +
		"from sys_menu m " +
		" left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"where m.status = '0' and rm.role_id = ?"
	return m.query(ctx, querySql, roleId)
}

func (m *menuRepository) query(ctx context.Context, querySql string, args ...any) ([]*model.SysMenu, error) {
	var menuList []*model.SysMenu
	result := m.DB(ctx).Raw(querySql, args...).Scan(&menuList)
	if result.Error != nil {
		return nil, result.Error
	}
	return menuList, nil
}
