package v1

var (
	// common errors
	ErrSuccess             = newError(200, "操作成功")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrForbid              = newError(403, "没有权限操作")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	// 签名相关
	ErrCheckSigFailed = newError(1001, "签名错误")
	ErrSigTimeExpired = newError(1002, "签名过期")
	ErrAesParamError  = newError(1003, "AES参数错误")
	ErrTokenError     = newError(1004, "认证错误")

	ErrParamError        = newError(1003, "参数错误")
	ErrLoginError        = newError(1004, "账号密码错误")
	ErrLoginExpired      = newError(1005, "登录过期")
	ErrUserNotFoundError = newError(1006, "用户不存在")
	ErrPwdError          = newError(1008, "密码错误")

	ErrEmailAlreadyUse    = newError(1005, "The email is already in use.")
	ErrUsernameAlreadyUse = newError(1006, "账号已存在")
	ErrCreateMenuFailed   = newError(1007, "创建菜单失败")
	ErrEditMenuFailed     = newError(1008, "更新菜单失败")
	ErrDelMenuFailed      = newError(1009, "删除菜单失败")
	ErrRolePermsFailed    = newError(1010, "获取角色失败")
)
