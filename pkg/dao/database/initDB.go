package database

import (
	"database/sql"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

//MysqlDriver mysql驱动器
const MysqlDriver = "mysql"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open(MysqlDriver, "root:123456@tcp(127.0.0.1:3306)/pixiu")
	// TODO: 从配置文件中获取数据库属性
	// var dsn string = conf.Configure["mysql_username"] + ":" + conf.Configure["mysql_password"] + "@tcp(" + conf.Configure["mysql_host"] + ":" + conf.Configure["mysql_port"] + ")/" + conf.Configure["mysql_dbname"] + "?charset=utf8"
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
}

func GetConnection() *sql.DB {
	return db
}
