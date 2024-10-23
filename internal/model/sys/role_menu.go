package model

type SysRoleMenu struct {
	/** 角色ID */
	RoleId int64 `json:"roleId"`
	/** 菜单ID */
	MenuId int64 `json:"menuId"`
}

func (rm *SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
