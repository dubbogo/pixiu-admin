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

	//// Get database properties from configuration file
	// username, password, host, port, dbname := config.Bootstrap.GetMysqlConfig()
	// panic(username)
	// dataSourceName := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8"
	// db, err = sql.Open(MysqlDriver, dataSourceName)
	db, err = sql.Open(MysqlDriver, "root:Fw900827@tcp(192.168.31.44:3306)/pixiu")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
}

func GetConnection() *sql.DB {
	return db
}
