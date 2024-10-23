package model

type SysRole struct {
	Id       int64  `json:"roleId"`
	RoleName string `json:"roleName"` //角色名称
	RoleKey  string `json:"roleKey"`  //角色权限
	RoleSort int    `json:"roleSort"` //角色排序
	//数据范围（1：所有数据权限；2：自定义数据权限；3：本部门数据权限；4：本部门及以下数据权限；5：仅本人数据权限）
	DataScope int `json:"dataScope"`
	Status    int `json:"status"`
	DelFlag   int `json:"delFlag"` //（0代表存在 2代表删除）
	//角色权限
	Permissions []string `json:"permissions" gorm:"-"`
}

func (r *SysRole) TableName() string {
	return "sys_role"
}
