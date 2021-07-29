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

func Register(c *gin.Context)  {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 格式检查， 前端实现？
	// 4到16位（字母，数字，下划线，减号）
	//if flag, _ := regexp.MatchString("^[a-zA-Z0-9_-]{4,16}$", username); !flag {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "illegal username",
	//		"data":   nil,
	//	})
	//	return
	//}
	// 密码至少包含 数字和英文，长度6-20
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
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    "注册失败: " + err.Error(),
		//	"data":    nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(err))
	}else {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": 0,
		//	"msg":    "success ",
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithRet("注册成功，请登录！"))
	}
}

// 登录结果
type LoginResult struct {
	Username string `json:"username"`
	Token string `json:"token"`
}


func Login(c *gin.Context)  {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//username := c.Query("username")
	//password := c.Query("password")

	password = utils.Md5(password)
	flag, _ := account.Login(username, password)
	if flag {
		generateToken(c, username)
	}else{
		//c.JSON(http.StatusOK, gin.H{
		//	"status": -1,
		//	"msg":    "验证失败, 登录信息有误",
		//	"data":   nil,
		//})
		c.JSON(http.StatusOK, controller.WithError(errors.New("验证失败, 登录信息有误!")))
	}
}

func generateToken(c *gin.Context, username string)  {
	j := auth.NewJWT()
	claims := auth.CustomClaims{
		username,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
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

	log.Println(token)

	data := LoginResult{
		Username: username,
		Token:    token,
	}
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 0,
	//	"msg":    "登录成功",
	//	"data":   data,
	//})
	c.JSON(http.StatusOK, controller.WithRet(data))
	return
}
