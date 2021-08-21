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
	"log"
	"net/http"
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller"
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller/auth"
	"github.com/dubbogo/pixiu-admin/pkg/logic/account"
	"github.com/dubbogo/pixiu-admin/pkg/utils"
)

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// username password Format check
	//if flag, _ := regexp.MatchString("^[a-zA-Z0-9_-]{4,16}$", username); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal username",
	//		"data":   nil,
	//	})
	//	return
	//}

	//if flag, _ := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", password); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal password",
	//		"data":   nil,
	//	})
	//	return
	//}
	password = utils.Md5(password)
	err := account.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, controller.WithError(err))
	} else {
		c.JSON(http.StatusOK, controller.WithRet("Register successfully, please login!"))
	}
}

// login result
type LoginResult struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//username := c.Query("username")
	//password := c.Query("password")

	password = utils.Md5(password)
	flag, _ := account.Login(username, password)
	if flag {
		generateToken(c, username)
	} else {
		c.JSON(http.StatusOK, controller.WithError(errors.New("Authentication failed, login information is wrong!")))
	}
}

func generateToken(c *gin.Context, username string) {
	j := auth.NewJWT()
	claims := auth.CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // Signature effective time
			ExpiresAt: int64(time.Now().Unix() + 3600), // Signature expiration time
			Issuer:    "dubbo-go-pixiu",
		},
	}
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, controller.WithError(err))
		return
	}

	log.Println(token)

	data := LoginResult{
		Username: username,
		Token:    token,
	}
	c.JSON(http.StatusOK, controller.WithRet(data))
}
