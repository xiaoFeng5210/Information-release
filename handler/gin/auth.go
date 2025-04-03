package handler

import (
	"infomation-release/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UID_IN_TOKEN = "uid"
	UID_IN_CTX   = "uid"
	COOKIE_NAME  = "jwt"
)

var (
	keyConfig = "123456"
)

func GetUidFromJwt(token string) int {
	_, payload, err := util.VerifyJwt(token, keyConfig)
	if err != nil {
		return 0
	}

	for k, v := range payload.UserDefined {
		if k == UID_IN_TOKEN {
			return int(v.(float64))
		}
	}

	return 0
}

func GetLoginUidFromCookie(ctx *gin.Context) int {
	token := ""
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == COOKIE_NAME {
			token = cookie.Value
		}
	}
	return GetUidFromJwt(token)
}

func Auth(ctx *gin.Context) {
	loginUid := GetLoginUidFromCookie(ctx)
	if loginUid <= 0 {
		ctx.Abort()
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "请先登录",
		})
		return
	} else {
		ctx.Set(UID_IN_CTX, loginUid)
		// ctx.Next() 可以省略
	}
}
