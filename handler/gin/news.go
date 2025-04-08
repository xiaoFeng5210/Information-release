package handler

import (
	database "infomation-release/database/gorm"
	"infomation-release/handler/model"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostNews(c *gin.Context) {
	loginId := c.Value(UID_IN_CTX).(int)

	var news model.News
	err := c.ShouldBindJSON(&news)

	if err != nil {
		slog.Error("PostNews failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1004,
			"msg":  "参数错误",
		})
		return
	}

	id, err := database.PostNews(loginId, news.Title, news.Content)
	if err != nil {
		slog.Error("PostNews failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1005,
			"msg":  "发布新闻失败",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "发布新闻成功",
		"data": id,
	})
}

func GetNewsByID(c *gin.Context) {
	idStr := c.Param("id")
	if id, err := strconv.Atoi(idStr); err != nil {
		slog.Error("GetNewsByID failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1004,
			"msg":  "参数错误,id不满足格式要求",
		})
		return
	} else {
		news := database.GetNewsByID(id)
		if news == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1005,
				"msg":  "新闻不存在",
			})
			return
		}
		user := database.GetUserByID(news.UserId)
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": -1005,
				"msg":  "用户不存在",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取新闻成功",
				"data": gin.H{
					"news": news,
					"user": user,
				},
			})
		}
	}
}

func DeleteNews(c *gin.Context) {
	loginId := c.Value(UID_IN_CTX).(int)
	idStr := c.Param("id")

	if id, err := strconv.Atoi(idStr); err != nil {
		slog.Error("DeleteNews failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1004,
			"msg":  "参数错误,id不满足格式要求",
		})
		return
	} else {
		if !newsBelongUser(id, loginId) {
			c.JSON(http.StatusForbidden, gin.H{
				"code": -1005,
				"msg":  "没有权限删除新闻",
			})
			return
		}
		err := database.DeleteNews(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1005,
				"msg":  "删除新闻失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "删除新闻成功",
		})
	}
}

func UpdateNews(c *gin.Context) {
	loginId := c.Value(UID_IN_CTX).(int)
	var news model.News
	err := c.ShouldBindJSON(&news)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1004,
			"msg":  "参数错误",
		})
		return
	}
	if !newsBelongUser(news.Id, loginId) {
		c.JSON(http.StatusForbidden, gin.H{
			"code": -1005,
			"msg":  "没有权限更新新闻",
		})
		return
	}
	err = database.UpdateNews(news.Id, news.Title, news.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1005,
			"msg":  "更新新闻失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新新闻成功",
	})
}

// 判断新闻nid是不是用户uid发布的
func newsBelongUser(nid, uid int) bool {
	news := database.GetNewsByID(nid)
	if news != nil {
		if news.UserId == uid {
			return true
		}
	}
	return false
}
