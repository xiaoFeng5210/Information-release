package handler

import (
	database "infomation-release/database/gorm"
	model "infomation-release/handler/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	// TODO:
	uid, ok := ctx.Value(UID_IN_CTX).(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "请先登陆",
		})
		return
	}
	var req model.ModifyPassRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	err := database.UpdatePassword(uid, req.NewPass, req.OldPass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码修改成功",
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
