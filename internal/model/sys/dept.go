package model

import "gorm.io/gorm"

type SysDept struct {
	gorm.Model
	Base
	DeptId     int64     `json:"dept_id"`
	ParentId   int64     `json:"parent_id"` //父部门ID
	ParentName string    `json:"parent_name"`
	DeptName   string    `json:"dept_name"`
	OrderNum   int       `json:"order_num"` //显示顺序
	Leader     string    `json:"leader"`    //负责人
	Phone      string    `json:"phone"`     //联系电话
	Email      string    `json:"email"`
	Status     int       `json:"status"`   //部门状态:0正常, 1停用
	Children   []SysDept `json:"children"` //子部门
}

func (d *SysDept) TableName() string {
	return "sys_dept"
}
