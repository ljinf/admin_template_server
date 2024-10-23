package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

//	type LoginRequest struct {
//		Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
//		Password string `json:"password" binding:"required" example:"123456"`
//	}

type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Id          int64  `json:"userId"`
	NickName    string `json:"nickName"` //用户昵称
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`    //新密码
	OldPassword string `json:"oldPassword"` //旧密码
	Sex         string `json:"sex"`         //性别  0=男,1=女,2=未知
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type EditUserRequest struct {
	Id          int64  `json:"userId"`
	DeptId      int64  `json:"deptId"`
	Username    string `json:"username" gorm:"username"` //登录名称
	NickName    string `json:"nickName"`                 //用户名称
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Sex         string `json:"sex"`     //性别  0=男,1=女,2=未知
	Avatar      string `json:"avatar"`  //头像
	Status      string `json:"status" ` //账号状态
	LoginIp     string `json:"-"`
	LoginDate   int64  `json:"-"` //最后登录时间戳

	RoleIds []int `json:"roleIds"`
}
