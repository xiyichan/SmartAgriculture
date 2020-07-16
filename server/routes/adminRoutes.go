package routes

import (
	"server/controller"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) *gin.Engine{
	r.POST("/api/admin/register",controller.AdminRegister)
	r.GET("/api/admin/confirm",controller.AdminConfirmRegister)


return r
}
