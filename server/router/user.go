package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
	"server/middleware"
)

func UserRouter(r *gin.Engine) *gin.Engine {
	r.POST("api/user/login/password", controllers.UserLoginByPassword)

	UserToken := r.Group("api/user").Use(middleware.TokenMiddleware())
	UserToken.POST("upload/avatar", controllers.UserUploadAvatar)
	UserToken.POST("update/password", controllers.UserUpdatePassword)
	UserToken.GET("info", controllers.UserInfo)
	UserToken.POST("update/name", controllers.UserUpdateUsername)

	UserEmailCaptcha := r.Group("api/user").Use(middleware.EmailCaptchaMiddleware())
	UserEmailCaptcha.POST("/register/email", controllers.UserRegisterByEmail)
	UserEmailCaptcha.POST("/forget/password/email", controllers.UserForgetPasswordByEmail)

	return r
}
