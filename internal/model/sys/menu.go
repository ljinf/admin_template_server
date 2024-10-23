package model

type SysMenu struct {
	//gorm.Model
	//Base
	MenuId    int64      `json:"menu_id"`           //菜单名称
	MenuName  string     `json:"menu_name"`         //菜单名称
	ParentId  int64      `json:"parent_id"`         //父菜单ID
	OrderNum  int        `json:"order_num"`         //排序
	Path      string     `json:"path"`              //路由地址
	Component string     `json:"component"`         //组件路径
	Query     string     `json:"query"`             //路由参数
	IsFrame   string     `json:"is_frame"`          //是否为外链（0是 1否）
	MenuType  string     `json:"menu_type"`         //类型（M目录 C菜单 F按钮）
	Visible   string     `json:"visible"`           //菜单状态（0显示 1隐藏）
	Status    string     `json:"status"`            //菜单状态（0正常 1停用）
	Perms     string     `json:"perms"`             //权限字符串
	Icon      string     `json:"icon"`              //菜单图标
	Children  []*SysMenu `json:"children" gorm:"-"` //子菜单
}

func (m *SysMenu) TableName() string {
	return "sys_menu"
}
