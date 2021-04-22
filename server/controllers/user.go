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

func UserRegisterByPhone(ctx *gin.Context) {

}
func UserRegisterByEmail(ctx *gin.Context) {
	db := common.GetDB()
	rdb := common.GetRedis()
	var requestUser = models.User{}
	ctx.Bind(&requestUser)
	if requestUser.Email != "" {
		if err := gvalid.Check(requestUser.Email, "email", nil); err != nil {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该邮箱格式不对"})
			return
		}
	}

	if err := gvalid.Check(requestUser.Password, "password", nil); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码格式不对,6位以上"})
		return
	}

	if len(requestUser.Name) == 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称为空"})
		return
	} else if len(requestUser.Name) > 20 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称过长"})
		return
	}

	var user models.User
	db.Where("email = ?", requestUser.Email).Find(&user)
	if user.ID > 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "已有该账号"})
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	newUser := models.User{
		Name:     requestUser.Name,
		Email:    requestUser.Email,
		Phone:    requestUser.Phone,
		Password: string(hashPassword),
	}
	if err := db.Create(&newUser).Error; err == nil {
		claims := common.Claims{
			ID:       newUser.ID,
			Email:    newUser.Email,
			Password: newUser.Password,
			Role:     "user",
			Name:     newUser.Name,
			Phone:    newUser.Phone,
		}
		if token, err := common.ReleaseToken(claims); err == nil {
			path := fmt.Sprintf("public/user/%d", newUser.ID)
			os.Mkdir(path, os.ModePerm)
			herr := rdb.HSet(ctx, string(newUser.ID), "token", token, "count", 0, "time", time.Now().Unix())
			if herr.Err() != nil {
				ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
				return
			}
			ctx.JSON(200, gin.H{"code": 200, "msg": "注册成功", "token": token})
		} else {
			ctx.JSON(500, gin.H{"code": 500, "msg": "token发放错误"})
		}
	} else {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误，数据库异常"})
		return
	}
}
func UserLoginByPassword(ctx *gin.Context) {
	db := common.GetDB()
	rdb := common.GetRedis()
	account := ctx.PostForm("Account")
	password := ctx.PostForm("Password")
	var user models.User
	//var dto models.UserDto
	if err := gvalid.Check(account, "email", nil); err == nil {
		user.Email = account
		db.Where("email=?", account).First(&user)
		if user.ID == 0 {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该账号不存在"})
			return
		}
		count, _ := rdb.HGet(ctx, "user"+strconv.Itoa(int(user.ID)), "count").Int()
		if count >= 5 {
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			//需要防止暴力破解
			if count >= 5 {
				//密码错误五次，冻结账号
				err := rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "count", count+1, "time", time.Now().Unix()+int64((count-4)*60)).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
				ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
				return
			} else {
				//记录密码错误次数
				err := rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "count", count+1, "time", time.Now().Unix()).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
			}
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
			return
		}

		claims := common.Claims{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Role:     "user",
			Phone:    user.Phone,
		}
		fmt.Println(claims)
		token, err := common.ReleaseToken(claims)
		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常,token生成失败"})
			return
		}

		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "msg": "系统异常，Redis异常"})
			return
		}
		//rdb.Del(ctx, "user"+strconv.Itoa(int(user.ID)))
		//TODO:测试
		rdb.HDel(ctx, "user"+strconv.Itoa(int(user.ID)), "count", "time")
		//fmt.Println(user.ID)
		a := models.ToUserDto(user)
		a.ID = user.ID
		//fmt.Println(a)
		//userDto:=models.UserDto{
		//	ID: user.ID,
		//
		//
		//}
		rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "token", token)
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": a})

		return
	}

	//电话登录

	if err := gvalid.Check(account, "phone", nil); err == nil {
		user.Phone = account
		db.Where("phone=?", account).First(&user)
		if user.ID == 0 {
			ctx.JSON(422, gin.H{"code": 422, "msg": "该账号不存在"})
			return
		}
		count, _ := rdb.HGet(ctx, "user"+strconv.Itoa(int(user.ID)), "count").Int()
		if count >= 5 {
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			//需要防止暴力破解
			if count >= 5 {
				//密码错误五次，冻结账号
				err := rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "count", count+1, "time", time.Now().Unix()+int64((count-4)*60)).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
				ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误5次，账号冻结，请使用找回密码重置密码"})
				return
			} else {
				//记录密码错误次数
				err := rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "count", count+1, "time", time.Now().Unix()).Err()
				if err != nil {
					ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误,Redis异常"})
					return
				}
			}
			ctx.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
			return
		}

		claims := common.Claims{
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Role:     "user",
			Phone:    user.Phone,
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
		//rdb.HDel(ctx, "user"+strconv.Itoa(int(user.ID)),"count","time")
		rdb.HDel(ctx, "user"+strconv.Itoa(int(user.ID)), "count", "time")
		rdb.HSet(ctx, "user"+strconv.Itoa(int(user.ID)), "token", token)
		ctx.JSON(200, gin.H{"code": 200, "msg": "登录成功", "token": token, "data": "models.ToUserDto(user)"})
		return
	}
	ctx.JSON(422, gin.H{"code": 422, "msg": "账号或者密码出错"})

}
func UserLoginByCaptcha(ctx *gin.Context) {

}
func UserUploadAvatar(ctx *gin.Context) {
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	//fmt.Println(reflect.TypeOf(claims))
	//fmt.Println(claims.(*common.Claims))
	c := claims.(*common.Claims)
	var u models.User
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

	path := fmt.Sprintf("public/user/%d/avatar/", c.ID)

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
func UserUpdatePassword(ctx *gin.Context) {
	rdb := common.GetRedis()
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	c := claims.(*common.Claims)
	newPassword := ctx.PostForm("Password")
	var u models.User
	u.ID = c.ID
	if err := gvalid.Check(newPassword, "password", nil); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码格式不对,6位以上"})
		return
	}
	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统密码处理错误"})
		return
	}
	if err := db.Model(&u).Update("password", hashNewPassword).Error; err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "修改失败"})
	} else {
		fmt.Println(u.ID)
		rdb.HDel(ctx, "user"+strconv.Itoa(int(c.ID)), "count", "time", "token")
		ctx.JSON(200, gin.H{"code": 200, "msg": "修改密码成功"})
	}

}
func UserForgetPasswordByPhone(ctx *gin.Context) {

}
func UserForgetPasswordByEmail(ctx *gin.Context) {
	rdb := common.GetRedis()
	db := common.GetDB()
	var user models.User
	email := ctx.PostForm("Email")
	password := ctx.PostForm("Password")

	db.Where("email=?", email).First(&user)
	if user.ID==0{
		ctx.JSON(422, gin.H{"code": 422, "msg": "没有该账号"})
		return
	}
	if err := gvalid.Check(password, "password", nil); err != nil {
		ctx.JSON(422, gin.H{"code": 422, "msg": "密码格式不对,6位以上"})
		return
	}
	hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统密码处理错误"})
		return
	}
	user.Password = string(hashNewPassword)

	if err := db.Model(&user).Update("password", user.Password).Error; err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "修改失败"})
	} else {
		rdb.HDel(ctx, "user"+strconv.Itoa(int(user.ID)), "count", "time", "token")
		//	rdb.HDel(ctx, "user"+strconv.Itoa(int(user.ID)),"count","time")
		ctx.JSON(200, gin.H{"code": 200, "msg": "修改密码成功"})
	}
}

func UserUpdateUsername(ctx *gin.Context) {
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	c := claims.(*common.Claims)
	newName := ctx.PostForm("Name")
	var u models.User
	u.ID = c.ID
	if len(newName) == 0 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称为空"})
	} else if len(newName) > 20 {
		ctx.JSON(422, gin.H{"code": 422, "msg": "昵称过长"})
	}
	if err := db.Model(&u).Update("name", newName).Error; err != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "修改失败"})
	}
	ctx.JSON(200, gin.H{"code": 200, "msg": "修改成功"})
}
func UserInfo(ctx *gin.Context) {
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	c := claims.(*common.Claims)
	var u models.User
	u.ID = c.ID
	db.First(&u)
	ctx.JSON(200, gin.H{"code": 200, "msg": "查找成功", "data": models.ToUserDto(u)})
}
func UserList(ctx *gin.Context) {
	claims, _ := ctx.Get("Claims")
	//fmt.Println(claims.(common.Claims).Role)
	if claims.(*common.Claims).Role != "admin" {
		ctx.JSON(422, gin.H{"code": 200, "msg": "权限不足"})
		return
	}
	db := common.GetDB()
	pageIndex := ctx.PostForm("PageIndex")
	pageSize := ctx.PostForm("PageSize")
	if pageIndex == "" || pageSize == "" {
		ctx.JSON(422, gin.H{"code": 422, "msg": "分页设置为空"})
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
		ctx.JSON(500, gin.H{"code": 500, "msg": "查询失败"})
	} else {
		ctx.JSON(200, gin.H{"code": 200, "data": dto, "msg": "查询成功", "count": count})
	}
}
