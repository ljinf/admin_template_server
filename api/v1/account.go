package v1

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"123456"`
	Password string `json:"password" binding:"required" example:"123456"`
	Code     string `json:"code"` //验证码
	Uuid     string `json:"uuid"` //唯一标识
}

type LoginResp struct {
	AccessToken string `json:"accessToken"`
}
