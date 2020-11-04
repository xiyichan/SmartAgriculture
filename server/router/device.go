package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
	"server/middleware"
)

func DeviceRouter(r *gin.Engine) *gin.Engine {



	pi := r.Group("api/device/pi").Use(middleware.TokenMiddleware())
	pi.POST("/register", controllers.RegisterPi)
	pi.POST("/set/property", controllers.UserSetPiProperty)
	pi.GET("/get/property", controllers.UserGetPiProperty)
	pi.POST("/bind/user", controllers.UserBindPi)
	pi.POST("user/list",controllers.UserPiList)
	pi.POST("list", controllers.PiList)
	pi.POST("delete", controllers.PiDelete)
	pi.POST("user/login/out",controllers.PiUserLoginOut)

	r.POST("api/device/pi/login/out",controllers.PiLoginOut)
	return r
}
