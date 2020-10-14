package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"server/common"
	"server/models"
	"time"
)

//登录

func AdminLoginByPassword(ctx *gin.Context) {
	db := common.GetDB()
	rdb := common.GetRedis()
	account := ctx.PostForm("Account")
	password := ctx.PostForm("Password")
	var admin models.Admin
	if err := gvalid.Check(account, "email", nil); err == nil {
		admin.Email = account
		db.Where("email=?", account).First(&admin)
		if admin.ID == 0 {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该账号不存在"})
			return
		}
		count, _ := rdb.HGet(ctx, "admin"+string(admin.ID), "count").Int()
		if count >= 5 {
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
			//需要防止暴力破解
			if count >= 5 {
				//密码错误五次，冻结账号
				err := rdb.HSet(ctx, "admin"+string(admin.ID), "count", count+1, "time", time.Now().Unix()+int64((count-4)*60)).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
				ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
				return
			} else {
				//记录密码错误次数
				err := rdb.HSet(ctx, "admin"+string(admin.ID), "count", count+1, "time", time.Now().Unix()).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
			}
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
			return
		}

		claims := common.Claims{
			ID:       admin.ID,
			Email:    admin.Email,
			Password: admin.Password,
			Name:     admin.Name,
			Role:     "admin",
			Phone:    admin.Phone,
		}
		token, err := common.ReleaseToken(claims)
		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常,token生成失败"})
			return
		}

		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常，Redis异常"})
			return
		}
		rdb.Del(ctx, "admin"+string(admin.ID))
		//rdb.HDel(ctx, "admin"+string(admin.ID),"count","time")
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": admin})
		return
	}

	if err := gvalid.Check(account, "phone", nil); err == nil {
		admin.Phone = account
		db.Where("phone=?", account).First(&admin)
		if admin.ID == 0 {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该账号不存在"})
			return
		}
		count, _ := rdb.HGet(ctx, "admin"+string(admin.ID), "count").Int()
		if count >= 5 {
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
			//需要防止暴力破解
			if count >= 5 {
				//密码错误五次，冻结账号
				err := rdb.HSet(ctx, "admin"+string(admin.ID), "count", count+1, "time", time.Now().Unix()+int64((count-4)*60)).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
				ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
				return
			} else {
				//记录密码错误次数
				err := rdb.HSet(ctx, "admin"+string(admin.ID), "count", count+1, "time", time.Now().Unix()).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
			}
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
			return
		}

		claims := common.Claims{
			Email:    admin.Email,
			Password: admin.Password,
			Name:     admin.Name,
			Role:     "admin",
			Phone:    admin.Phone,
		}
		token, err := common.ReleaseToken(claims)
		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常,token生成失败"})
			return
		}

		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常，Redis异常"})
			return
		}
		//rdb.HDel(ctx, "admin"+string(admin.ID),"count","time")
		rdb.Del(ctx, "admin"+string(admin.ID))
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": admin})
		return
	}
	ctx.JSON(422, gin.H{"code": 200, "msg": "账号或者密码出错"})

}

//添加管理员，从管理员页面新添一个。
func AdminAddByEmail(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	fmt.Println(claims.(common.Claims).Role)
	if claims.(common.Claims).Role != "admin" {
		ctx.JSON(422, gin.H{"code": 200, "msg": "权限不足"})
		return
	}
	db := common.GetDB()
	var requestAdmin = models.Admin{}
	ctx.Bind(&requestAdmin)
	if requestAdmin.Email != "" {
		if err := gvalid.Check(requestAdmin.Email, "email", nil); err != nil {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该邮箱格式不对"})
			return
		}
	}

	if err := gvalid.Check(requestAdmin.Password, "password", nil); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码格式不对,6位以上"})
		return
	}

	if len(requestAdmin.Name) == 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称为空"})
		return
	} else if len(requestAdmin.Name) > 20 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称过长"})
		return
	}

	var admin models.Admin
	db.Where("email = ?", requestAdmin.Email).Find(&admin)
	if admin.ID > 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "已有该账号"})
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(requestAdmin.Password), bcrypt.DefaultCost)
	newAdmin := models.Admin{
		Name:     requestAdmin.Name,
		Email:    requestAdmin.Email,
		Phone:    requestAdmin.Phone,
		Password: string(hashPassword),
	}
	if err := db.Create(&newAdmin).Error; err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "数据库出错"})

	}
	path := fmt.Sprintf("public/admin/%d", newAdmin.ID)
	os.Mkdir(path, os.ModePerm)
	ctx.JSON(200, gin.H{"code": 200, "msg": "添加"})
}

//修改个人信息

//
