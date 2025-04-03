package handler

import (
	database "infomation-release/database/gorm"
	model "infomation-release/handler/model"
	"infomation-release/util"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	COOKIE_LIFE = 7 * 86400
)

func RegisterUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	err := database.RegisterUser(user.Name, user.PassWord)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "注册成功",
	})
}

func UpdatePassword(ctx *gin.Context) {
	uid, ok := ctx.Value(UID_IN_CTX).(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请先登录",
		})
		return
	}

	var req model.ModifyPassRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1001,
			"msg":  "参数错误",
		})
		return
	}
	err = database.UpdatePassword(uid, req.NewPass, req.OldPass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1005,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码修改成功",
	})
}

func Login(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1001,
			"msg":  "参数错误",
		})
		return
	}

	user2 := database.GetUserByUsername(user.Name)
	if user2 == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1002,
			"msg":  "用户不存在",
		})
		return
	}
	if user2.PassWord != user.PassWord {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1003,
			"msg":  "密码错误",
		})
		return
	}

	// 接下来就是生成token并且种到cookie里
	header := util.DefaultHeader
	payload := util.JwtPayload{ //payload以明文形式编码在token中，server用自己的密钥可以校验该信息是否被篡改过
		Issue:       "news",
		IssueAt:     time.Now().Unix(),                                //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(COOKIE_LIFE * time.Second).Unix(), //7天后过期
		UserDefined: map[string]any{UID_IN_TOKEN: user2.Id},           //用户自定义字段。如果token里包含敏感信息，请结合https使用
	}

	if token, err := util.GenJWT(header, payload, keyConfig); err != nil {
		slog.Error("生成token失败", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1005,
			"msg":  "生成token失败",
		})
		return
	} else {
		ctx.SetCookie(
			COOKIE_NAME,
			token,       //注意：受cookie本身的限制，这里的token不能超过4K
			COOKIE_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
			"/",         //path，cookie存放目录
			"localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
			false,       //是否只能通过https访问
			true,        //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
		)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "登录成功",
		})
	}
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie(COOKIE_NAME, "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "退出登陆成功",
	})
}

func GetAllUsers(ctx *gin.Context) {
	// 返回指针，可以自动解引用和序列化
	users := database.GetAllUsers()
	if users != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "获取用户成功",
			"data": users,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "获取用户失败",
		})
	}
}

func GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "id不能为空",
		})
		return
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "id格式错误",
		})
		return
	}
	user := database.GetUserByID(uid)
	if user != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	}
}

func GetUserInfo(ctx *gin.Context) {
	loginUid := GetLoginUidFromCookie(ctx)
	if loginUid > 0 {
		user := database.GetUserByID(loginUid)
		if user != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取用户信息成功",
				"user": user,
			})
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": -1001,
			"msg":  "未找到用户，请先登陆或者查看用户是否存在",
		})
	}
}
