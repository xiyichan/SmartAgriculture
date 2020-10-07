package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"te/common"
	"time"
)

//TODO:该问题是如果密码修改之后原本的token依旧能用，若异地登录的话一样改密码也防止不了，使用黑名单或者从数据库查找密码同样都是需要查找数据库。
func TokenMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		if claims.ExpiresAt-time.Now().Unix() <= 5*60*60 {
			//在token过期前五个小时，就让其软失效
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		if len(roles) > 0 {
			allow := false
			for _, role := range roles {
				if claims.Role == role {
					allow = true
				}
			}
			if claims.Role == "admin" {
				allow = true
			}
			if allow == false {
				ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
				ctx.Abort()
				return
			}
		}

		ctx.Set("Claims", claims)
		ctx.Next()
	}
}
