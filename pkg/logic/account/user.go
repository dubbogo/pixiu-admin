package account


func EditPassword(oldPassword, newPassword, username string) (bool, error) {
	return userDao.EditPassword(oldPassword,newPassword,username)
}

func GetUserInfo(username string) (bool, interface{}, error) {
	return userDao.GetUserInfo(username)
}

func GetUserRole(username string) (bool, interface{}, error) {
	return userDao.GetUserRole(username)
}

func CheckUserIsAdmin(username string) (bool, error)  {
	return userDao.CheckUserIsAdmin(username)
}