package controller

import (
	"server/model"
	"server/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
	"golang.org/x/crypto/bcrypt"
	"os"
	path2 "path"
	"time"
)

func AdminRegister(context *gin.Context){
	DB:=common.GetDB()
	var requestAdmin=model.Admin{}
	context.Bind(&requestAdmin)
	if requestAdmin.Email!="" {
		if err := gvalid.Check(requestAdmin.Email, "email", nil); err != nil {
			context.JSON(422,gin.H{"code":422,"msg":"该邮箱格式不对"})
			return
		}
	} else if requestAdmin.Phone!=""{
		if err := gvalid.Check(requestAdmin.Phone, "phone", nil); err != nil {
			context.JSON(422,gin.H{"code":422,"msg":"该手机号码格式不对"})
			return
		}
	}
	if len(requestAdmin.NickName)==0{
		context.JSON(422,gin.H{"code":422,"msg":"昵称为空"})
		return
	} else if len(requestAdmin.NickName)>20{
		context.JSON(422,gin.H{"code":422,"msg":"昵称过长"})
		return
	}
	var admin model.Admin
	//判断数据库是否有该账号
	DB.Where("email = ?", requestAdmin.Email).Or("phone = ?", requestAdmin.Phone).Find(&admin)
	if admin.Uuid!=""{
		context.JSON(422,gin.H{"code":422,"msg":"已有该账号"})
		return
	}
	
	//发送邮件
	common.SendConfirmAdmin(requestAdmin)
	
	context.JSON(200,gin.H{"code":200,"msg":"等待审核，会以邮件邮件通知"})

}
func AdminConfirmRegister(context *gin.Context){
	email:=context.PostForm("Email")
	context.JSON(200,gin.H{"code":200,"msg":email})

}
func AdminLogin(context *gin.Context){
	DB:=common.GetDB()
	rdb,ctx:=common.GetRedis()
	var requestAdmin=model.Admin{}
	context.Bind(&requestAdmin)
	var admin model.Admin
	DB.Where("email=?",requestAdmin.Email).Or("phone=?",requestAdmin.Phone).First(&admin)


	if admin.Uuid==""{
		context.JSON(422,gin.H{"code":422,"msg":"不存在该账号"})
		return
	}

	//判断是否冻结了
	t,err:=rdb.HGet(ctx,admin.Uuid,"time").Int64()
//	fmt.Println(err)
	if err!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		panic(err)
		return
	}else if t>time.Now().Unix(){
		context.JSON(400, gin.H{"code": 400, "msg": "冻结","time":t})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(requestAdmin.Password)) ; err != nil {
		context.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
		//需要防止暴力破解

		count,err:=rdb.HGet(ctx,admin.Uuid,"count").Int()
		if err!=nil{
			panic(err)
		}
		//fmt.Println(count)
		if count>=5{
			hset:=rdb.HSet(ctx,admin.Uuid,"count",count+1,"time",time.Now().Unix()+int64((count-4)*60))
			panic(hset.Err())
		}else{
			hset:=rdb.HSet(ctx,admin.Uuid,"count",count+1,"time",time.Now().Unix())
			panic(hset.Err())
		}
		return
	}
	token,err:=common.ReleaseAdminToken(admin)
	if err!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		return
	}
	hset:=rdb.HSet(ctx,admin.Uuid,"token",token,"count",0,"time",time.Now().Unix())
	if hset.Err()!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		panic(hset.Err())
		return
	}
	//需要处理返回什么信息
	context.JSON(200,gin.H{"code":200,"msg":"登录成功","token":token,"data":model.ToAdminDto(admin)})
}
func AdminUploadAvater(context *gin.Context){
	DB:=common.GetDB()
	admin,_:=context.Get("admin")
	u:=admin.(model.Admin)
	avater,_:=context.FormFile("Avater")

	ext:=path2.Ext(avater.Filename)
	if ext!=".jpg"&&ext!=".jpeg"&&ext!=".bmp"&&ext!=".png"{
		context.JSON(400,gin.H{"code":400,"msg":"图片格式出错"})
		return
	}
	//判断图片大小 需要小于500K
	size:=avater.Size
	if size>512000{
		context.JSON(400,gin.H{"code":400,"msg":"图片过大"})
		return
	}
	path:=fmt.Sprintf("public/admin/%s/avater/",u.Uuid)
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	if err := context.SaveUploadedFile(avater, path+avater.Filename); err != nil {
		context.JSON(400,gin.H{"code":400,"msg":"上传出错"})
		return
	}
	err:=DB.Model(&u).Update("avater",path+avater.Filename)
	if err!=nil{
		context.JSON(200,gin.H{"code":200,"msg":"上传成功"})
	}else{
		context.JSON(500,gin.H{"code":500,"msg":"系统出错"})
	}
}
func AdminUpdatePassword(context *gin.Context){

}
func AdminForgetPassWord(context *gin.Context){

}
func AdminInfo(context *gin.Context){
	admin,_:=context.Get("admin")
	context.JSON(200,gin.H{"code":200,"msg":"获取成功","data":model.ToAdminDto(admin.(model.Admin))})
}
func AdminEmailVerify(context *gin.Context){
	email:=context.PostForm("Email")
	rdb,ctx:=common.GetRedis()

	captcha,err:=common.SendCaptcha(email)

	if err!=nil{
		context.JSON(500,gin.H{"code":500,"msg":"系统错误"})
	}

	hset:=rdb.Set(ctx,email,captcha,time.Minute*2)
	if hset.Err()!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		panic(hset.Err())
		return
	}

	context.JSON(200,gin.H{"code":200,"msg":"发送成功"})
}
func AdminPhoneVerify(context *gin.Context){

}
