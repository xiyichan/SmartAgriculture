package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
)

func UtilRouter(r *gin.Engine) *gin.Engine {
	r.POST("/api/email/captcha", controllers.SendEmailCaptcha)
	return r
}
