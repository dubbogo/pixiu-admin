package account

import (
	"github.com/dubbogo/pixiu-admin/pkg/dao"
	"github.com/dubbogo/pixiu-admin/pkg/dao/impl"
)

var(
	guestDao dao.GuestDao = impl.NewGuestDao()
	userDao dao.UserDao = impl.NewUserDao()
)

func Login(username string, password string) (bool, int) {
	return guestDao.Login(username, password)
}

func Register(username, password string) error {
	return guestDao.Register(username, password)
}
