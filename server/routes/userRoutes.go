package routes

import (
	"server/controller"
	"github.com/gin-gonic/gin"
	"server/middleware"
)

func UserRoutes(r *gin.Engine) *gin.Engine{
	//r.POST("/api/user/register", middleware.UserCaptchaMiddleware(),controller.UserRegister)

	r.POST("/api/user/email/verify",controller.UserEmailVerify)
	r.POST("/api/user/login",controller.UserLogin)

	userToken:=r.Group("api/user")
	userToken.Use(middleware.UserTokenMiddleware())
	userToken.GET("/info",controller.UserInfo)
	userToken.POST("/upload/avater",controller.UserUploadAvater)
	userToken.POST("/update/password",controller.UserUpdatePassword)

	userVerify:=r.Group("api/user")
	userVerify.Use(middleware.UserCaptchaMiddleware())
	userVerify.POST("/register",controller.UserRegister)
	userVerify.POST("/forgetpassword",controller.UserForgetPassword)



	return r
}
