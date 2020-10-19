package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path"
	"server/common"
	"server/models"
	"strconv"
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
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": models.ToAdminDto(admin)})
		return
	}

	//电话登录
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
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": models.ToAdminDto(admin)})
		return
	}
	ctx.JSON(422, gin.H{"code": 200, "msg": "账号或者密码出错"})

}

//添加管理员，从管理员页面新添一个。
//
func AdminAddByEmail(ctx *gin.Context) {
	claims, _ := ctx.Get("Claims")
	//fmt.Println(claims.(common.Claims).Role)
	if claims.(*common.Claims).Role != "admin" {
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


func AdminUploadAvatar(ctx *gin.Context){
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	if claims.(*common.Claims).Role != "admin" {
		ctx.JSON(422, gin.H{"code": 200, "msg": "权限不足"})
		return
	}
	c := claims.(*common.Claims)
	var u models.Admin
	u.ID = c.ID
	avatar, _ := ctx.FormFile("Avatar")
	ext := path.Ext(avatar.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".bmp" && ext != ".png" {
		ctx.JSON(400, gin.H{"code": 400, "msg": "图片格式出错"})
		return
	}
	size := avatar.Size
	if size > 512000 {
		ctx.JSON(400, gin.H{"code": 400, "msg": "图片过大"})
		return
	}

	path := fmt.Sprintf("public/admin/%d/avatar/", c.ID)

	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	if err := ctx.SaveUploadedFile(avatar, path+avatar.Filename); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "msg": "上传出错"})
		return
	}
	err := db.Model(&u).Update("avatar", path+avatar.Filename)
	if err != nil {
		ctx.JSON(200, gin.H{"code": 200, "msg": "上传成功"})
	} else {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统出错"})
	}
}

func AdminUpdatePassword(ctx *gin.Context){

}

func UserList(context *gin.Context) {
	db := common.GetDB()
	pageIndex := context.PostForm("PageIndex")
	pageSize := context.PostForm("PageSize")
	if pageIndex == "" || pageSize == "" {
		context.JSON(422, gin.H{"code": 422, "msg": "分页设置为空"})
		return
	}
	pi, _ := strconv.Atoi(pageIndex)
	ps, _ := strconv.Atoi(pageSize)

	//var users []models.User
	var dto []models.UserDto

	var count int64
	//正序
	db.Model(&models.User{}).Count(&count)
	e := db.Model(&models.User{}).Offset((pi - 1) * ps).Limit(ps).Find(&dto).Error
	if e != nil {
		context.JSON(500, gin.H{"code": 500, "msg": "查询失败"})
	} else {
		context.JSON(200, gin.H{"code": 200, "data": dto, "msg": "查询成功","count":count})
	}
}

func PiList(context *gin.Context){
	db := common.GetDB()
	pageIndex := context.PostForm("PageIndex")
	pageSize := context.PostForm("PageSize")
	if pageIndex == "" || pageSize == "" {
		context.JSON(422, gin.H{"code": 422, "msg": "分页设置为空"})
		return
	}
	pi, _ := strconv.Atoi(pageIndex)
	ps, _ := strconv.Atoi(pageSize)
	var dto[]models.PiListData

	var count int64
	db.Model(&models.Pi{}).Count(&count)
	e := db.Model(&models.Pi{}).Offset((pi - 1) * ps).Limit(ps).Find(&dto).Error
	if e != nil {
		context.JSON(500, gin.H{"code": 500, "msg": "查询失败"})
	} else {
		context.JSON(200, gin.H{"code": 200, "data": dto, "msg": "查询成功","count":count})
	}
}
