package router

import (
	"github.com/gin-gonic/gin"
	"server/controllers"
)

func WebSocketRouter(r *gin.Engine) *gin.Engine {

	r.GET("/ws", controllers.WSHandler)
	return r
}
