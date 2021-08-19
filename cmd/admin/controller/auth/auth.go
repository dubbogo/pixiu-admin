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

package auth

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
)

// 检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			//c.JSON(http.StatusOK, gin.H{
			//	"status": -1,
			//	"msg":    "请求未携带token, 无权限访问",
			//	"data":   nil,
			//})
			c.JSON(http.StatusOK, controller.WithError(errors.New("请求未携带token, 无权限访问")))
			c.Abort()
			return
		}
		log.Print("get token: ", token)
		j := NewJWT()
		// 解析token中包含信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// token 授权过期情况
			if err == TokenExpired {
				//c.JSON(http.StatusOK, gin.H{
				//	"status": -1,
				//	"msg": 	  "token授权已过期, 请重新申请授权",
				//	"data":   nil,
				//})
				c.JSON(http.StatusOK, controller.WithError(errors.New("token授权已过期, 请重新申请授权")))
				c.Abort()
				return
			}
			// 其他token错误情况
			//c.JSON(http.StatusOK, gin.H{
			//	"status": -1,
			//	"msg":    err.Error(),
			//	"data":   nil,
			//})
			c.JSON(http.StatusOK, controller.WithError(err))
		}
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token is not valid yet")
	TokenMalformed   error  = errors.New("This is not a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token")
	SignKey          string = "dubbo-go-pixiu" // TODO: 签名信息设置为动态获取
)

// 自定义载荷
type CustomClaims struct {
	Username string `json:"username"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// 新建jwt示例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取SignKey
func GetSignKey() string {
	return SignKey
}

// CreateToken 生成token(基于用户基本信息)
// 采用HS256算法
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// 返回token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 输入:token字符串, 自定义的Claims结构体对象,自定义函数
	// 解析token字符串为jwt的Token结构体指针
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	// 将token中的claims信息解析出来和用户原始数据进行校验, 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 过期时间验证
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// 设置token过期时间
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
