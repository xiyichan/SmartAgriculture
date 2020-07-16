package middleware

import (
	"server/model"
	"server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserTokenMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 获取authorization header
		rds, ctx := common.GetRedis()
		tokenString := context.GetHeader("Authorization")
		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseUserToken(tokenString)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}
		// 验证通过后获取claim 中的uuid
		u1 := claims.Uuid
		DB := common.GetDB()
		var user model.User
		DB.Where("uuid=?", u1).First(&user)

		// 用户
		if user.Uuid == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			context.Abort()
			return
		}
		rdsToken, err := rds.HGet(ctx, user.Uuid, "token").Result()
		if err != nil {
			context.JSON(500, gin.H{"code": 401, "msg": "系统出错"})
		}
		if rdsToken != tokenString {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "已在别的地方登录"})
			context.Abort()
			return
		}
		// 用户存在 将user 的信息写入上下文
		context.Set("user", user)

		context.Next()
	}
}

func UserCaptchaMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		DB:=common.GetDB()
		rds, ctx := common.GetRedis()
		captcha := context.PostForm("Captcha")
		email := context.PostForm("Email")
		phone := context.PostForm("Phone")
		var user model.User
		//通过邮箱验证
		c1, err := rds.Get(ctx, email).Result()
		if err == nil {
			if c1 == captcha {
				rds.Del(ctx, email)
				DB.Where("email=?",email).First(&user)
				context.Set("user", user)
				context.Next()
				return
			}
		} else if err.Error() == "redis: nil" {
			//注意这个错误是没有该值，中间有个空格
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "验证不通过"})
			context.Abort()
			return
		} else {
			context.JSON(500, gin.H{"code": 500, "msg": "系统出错"})
			context.Abort()
			return
		}
		//通过手机验证
		c2, err := rds.Get(ctx, phone).Result()
		if err == nil {
			if c2 == captcha {
				rds.Del(ctx, phone)
				DB.Where("phone",phone).First(&user)
				context.Set("user", user)
				context.Next()
				return
			}
		} else if err.Error() == "redis: nil" {
			//注意这个错误是没有该值，中间有个空格
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "验证不通过"})
			context.Abort()
			return
		} else {
			context.JSON(500, gin.H{"code": 500, "msg": "系统出错"})
			context.Abort()
			return
		}
		context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "验证不通过"})
		context.Abort()

	}
}
