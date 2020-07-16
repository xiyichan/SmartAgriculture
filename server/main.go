package main

import (

	"server/common"
	"server/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"net/http"
	"os"
)
var r *gin.Engine
func main() {
	InitConfig()
	InitApplication()
	InitRoutes()

}

func InitConfig()  {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func InitApplication(){
	common.InitDB()
	common.InitSmtp()
	common.InitRedis()
	common.InitIot()

}
func InitRoutes(){
	r = gin.Default()
	r=routes.UserRoutes(r)
	r=routes.AdminRoutes(r)
	r.StaticFS("/public",http.Dir("public"))
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run	())
}