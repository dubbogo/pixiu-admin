package dao

// GuestDao
type GuestDao interface {
	// Login 登录
	Login(username, password string) (bool, int)
	// CheckLogin 检查用户是否登录
	CheckLogin()
	// Register 用户注册
	Register(username, password string) error
}


// UserDao
type UserDao interface {
	// EditPassword 修改用户密码
	EditPassword(oldPassword, newPassword, username string) (bool, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(username string) (bool, interface{}, error)
	// GetUserRole 获取用户角色
	GetUserRole(username string) (bool, interface{}, error)
	// CheckUserIsAdmin 判断用户是否管理员
	CheckUserIsAdmin(username string) (bool,error)
}
