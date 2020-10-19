package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
	"server/middleware"
)

func AdminRouter(r *gin.Engine) *gin.Engine {
	r.POST("api/admin/login/password", controllers.AdminLoginByPassword)
	//r.POST("api/admin/register",controllers.AdminAddByEmail).Use(middleware.CorsMiddleware())

	adminToken:=r.Group("api/admin").Use(middleware.TokenMiddleware())
	adminToken.POST("register",controllers.AdminAddByEmail,middleware.EmailCaptchaMiddleware())
	adminToken.POST("user/list",controllers.UserList)
	adminToken.POST("upload/avatar",controllers.AdminUploadAvatar)
	adminToken.POST("pi/list",controllers.PiList)
	return r
}
