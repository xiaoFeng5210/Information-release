package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	if err := ctx.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "文件上传成功",
	})
}

func DownloadFile(ctx *gin.Context) {
	filename := ctx.Query("filename")
	fmt.Println(filename)
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	// 设置响应头，告诉浏览器这是一个文件下载
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/octet-stream")

	ctx.File("./uploads/" + filename)
}
