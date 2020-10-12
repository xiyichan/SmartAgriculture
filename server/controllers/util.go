package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
	"server/common"
	"time"
)

func SendEmailCaptcha(ctx *gin.Context) {
	rdb := common.GetRedis()
	email := ctx.PostForm("Email")
	if err := gvalid.Check(email, "email", nil); err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "邮箱格式有误"})
		return
	}
	captcha, err := common.SendCaptcha(email)
	if err != nil {
		//fmt.Println(err.Error())
		ctx.JSON(500, gin.H{"code": 500, "msg": "邮件发送失败"})
		return
	}
	//验证码存入redis，五分钟后过期
	if err = rdb.Set(ctx, email, captcha, time.Minute*5).Err(); err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误，Redis异常"})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "验证码发送成功，五分钟后过期"})

}
func RefreshToken(ctx *gin.Context) {

}
