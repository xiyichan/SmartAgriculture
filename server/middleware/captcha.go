package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/common"
)

func EmailCaptchaMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rds := common.GetRedis()
		captcha := ctx.PostForm("Captcha")
		email := ctx.PostForm("Email")
		//通过邮箱验证，验证码经由UserSendEmail发送，存在redis中
		c1, err := rds.Get(ctx, email).Result()
		if err == nil {
			if c1 == captcha {
				rds.Del(ctx, email)
				ctx.Next()
				return
			}
		} else if err.Error() == "redis: nil" {
			//注意这个错误是没有该值，中间有个空格
			//而不是没有email这个key
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "验证不通过"})
			ctx.Abort()
			return
		} else {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统出错"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "验证不通过"})
		ctx.Abort()
	}

}
