package model

type SysRoleUser struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
}

func (r *SysRoleUser) TableName() string {
	return "sys_user_role"
}
