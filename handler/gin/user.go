package handler

import (
	database "infomation-release/database/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
