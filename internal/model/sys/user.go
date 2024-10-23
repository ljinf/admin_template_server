package model

import "time"

type SysUser struct {
	//Base
	Id          int64     `json:"userId"`
	DeptId      int64     `json:"deptId"`
	Username    string    `json:"username" gorm:"username"` //登录名称
	NickName    string    `json:"nickName"`                 //用户名称
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Password    string    `json:"password"`
	Sex         string    `json:"sex"`     //性别  0=男,1=女,2=未知
	Avatar      string    `json:"avatar"`  //头像
	Status      string    `json:"status" ` //账号状态
	LoginIp     string    `json:"loginIp"`
	LoginDate   time.Time `json:"loginDate"` //最后登录时间戳
}

func (u *SysUser) TableName() string {
	return "sys_user"
}
