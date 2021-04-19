package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/common"
	"strconv"
	"strings"
	"time"
)

//TODO:该问题是如果密码修改之后原本的token依旧能用，若异地登录的话一样改密码也防止不了，使用黑名单或者从数据库查找密码同样都是需要查找数据库。
func TokenMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rdb := common.GetRedis()
		tokenString := ctx.GetHeader("Authorization")
		fmt.Println("----------------------", tokenString)
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
		rDBToken, err := rdb.HGet(ctx, "user"+strconv.Itoa(int(claims.ID)), "token").Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "系统出错"})
			return
		}
		if rDBToken != tokenString {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "已在别的地方登录"})
			ctx.Abort()
			return
		}
		ctx.Set("Claims", claims)
		fmt.Println(claims)
		ctx.Next()
	}
}
