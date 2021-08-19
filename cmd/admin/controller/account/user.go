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

package account

import (
	"net/http"
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

import (
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller"
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller/auth"
	"github.com/dubbogo/pixiu-admin/pkg/logic/account"
)

// Logout 用户注销
func Logout(c *gin.Context){
	// 设置token无效
	j := auth.NewJWT()
	claims := auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix()), // 签名过期时间
			Issuer:    "dubbo-go-pixiu",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    err.Error(),
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}
	// TODO 优化返回的json
	//c.JSON(http.StatusOK, gin.H{
	//	"status": -1,
	//	"msg":    "logout",
	//	"token":   token,
	//})
	c.JSON(http.StatusOK, controller.WithRet(token))
}


// EditPassword 修改账户密码
func EditPassword(c *gin.Context) {

	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")

	//if flag, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", oldPassword); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal oldPassword",
	//		"data":   nil,
	//	})
	//	return
	//}
	//
	//if flag, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", newPassword); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal newPassword",
	//		"data":   nil,
	//	})
	//	return
	//}
	username := c.Request.Header.Get("username")
	flag, err := account.EditPassword(oldPassword, newPassword, username)
	if !flag {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    err.Error(),
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}
	c.JSON(http.StatusOK, controller.WithRet("修改密码成功！"))
	// TODO 是否需要更新token？
	//generateToken(c, username)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	username := c.Request.Header.Get("username")
	flag, userInfo, err := account.GetUserInfo(username)
	if !flag {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    err.Error(),
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 0,
	//	"msg":    "获取用户信息成功",
	//	"data":   userInfo,
	//})
	c.JSON(http.StatusOK, controller.WithRet(userInfo))
}

// GetUserRole 获取用户角色
func GetUserRole(c *gin.Context)  {
	username := c.Request.Header.Get("username")
	flag, result, err := account.GetUserRole(username)
	if !flag {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    err.Error(),
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 0,
	//	"msg":    "获取用户角色成功",
	//	"data":   result,
	//})
	c.JSON(http.StatusOK, controller.WithRet(result))
}

// CheckUserIsAdmin 判断是否为管理员
func CheckUserIsAdmin(c *gin.Context) {
	username := c.Request.Header.Get("username")
	flag, err := account.CheckUserIsAdmin(username)
	if !flag {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    err.Error(),
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 0,
	//	"msg":    "This user is admin",
	//	"data":   nil,
	//})
	c.JSON(http.StatusOK, controller.WithRet("This user is admin"))
}
