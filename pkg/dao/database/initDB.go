/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"database/sql"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// MysqlDriver
const MysqlDriver = "mysql"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open(MysqlDriver, "root:123456@tcp(127.0.0.1:3306)/pixiu")
	// TODO: Get database properties from configuration file
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
