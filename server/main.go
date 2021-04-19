package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"server/common"
	"server/middleware"
	"server/router"
)

var r *gin.Engine

func init() {
	workDir, _ := os.Getwd()
	//viper 配置处理框架
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	common.InitMysql()
	common.InitRedis()
	common.InitSmtp()
	//	common.InitMongodb()
	common.InitAliyunIot()
	go common.InitAmqp()
	// 记录到文件。

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r = gin.Default()
	r.Use(middleware.CorsMiddleware())
	r = router.UtilRouter(r)
	r = router.UserRouter(r)
	r = router.WebSocketRouter(r)
	r = router.AdminRouter(r)
	r = router.DeviceRouter(r)
	r.StaticFS("/public", http.Dir("public"))

}

func main() {

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
	r.Run(":" + port)

}
