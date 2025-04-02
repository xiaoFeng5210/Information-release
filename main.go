package main

import (
	database "infomation-release/database/gorm"
	handler "infomation-release/handler/gin"
	"infomation-release/util"

	"github.com/gin-gonic/gin"
)

func main() {
	util.InitSlog("./log/ir.log")
	database.ConnectDB("./conf", "db", "yaml", "./log")
	engine := gin.Default()

	engine.GET("/download", handler.DownloadFile)

	engine.GET("/users", handler.GetAllUsers)

	if err := engine.Run(":3154"); err != nil {
		panic(err)
	}

}
