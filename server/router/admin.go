package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
)

func AdminRouter(r *gin.Engine) *gin.Engine {
	r.POST("api/admin/login/password", controllers.AdminLoginByPassword)
	r.POST("api/admin/register")
	return r
}
