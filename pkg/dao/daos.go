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
	CheckUserIsAdmin(username string) (bool, error)
}
