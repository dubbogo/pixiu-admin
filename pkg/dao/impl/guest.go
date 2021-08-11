package impl

import (
	SQL "database/sql"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/dao"
	"github.com/dubbogo/pixiu-admin/pkg/dao/database"
)

type GuestDao struct {
	db *SQL.DB
}

// TODO
func NewGuestDao() *GuestDao {
	return &GuestDao{
		db: database.GetConnection(),
	}
}

func (d *GuestDao)Create(db *SQL.DB) (interface{}, error){
	d.db = db
	var i dao.GuestDao = d
	return &i, nil
}


func (d *GuestDao) Login(username, password string) (bool, int){
	db := d.db
	var id int
	err := db.QueryRow("SELECT id FROM pixiu_user WHERE username = ? AND password = ?;", username, password).Scan(&id)
	if err != nil {
		return false, 0
	}
	return true, id
}

func (d *GuestDao) Register(username, password string) error {

	if username == "" {
		return errors.New("void username")
	}

	db := d.db
	var id int
	err := db.QueryRow("SELECT id FROM pixiu_user WHERE username = ?", username).Scan(&id)
	if err == nil {
		return errors.New("用户已存在, 请登录")
	}
	//now := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO pixiu_user (username,password) VALUES (?,?)")
	if err != nil {
		return errors.New("Illegal SQL statement!")
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password)
	if err != nil {
		return errors.New("Failed to create data!")
	}
	// TODO 设置事务， 动态设置用户角色
	err = db.QueryRow("SELECT id FROM pixiu_user WHERE username = ?", username).Scan(&id)
	stmt, err = db.Prepare("INSERT INTO pixiu_user_role(user_id, role_id) VALUE (?, ?)")
	_, err = stmt.Exec(id, 1)
	return err
}


func (d *GuestDao) CheckLogin() {
	panic("implement me")
}
