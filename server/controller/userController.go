package controller

import (
	"server/common"
	"server/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gvalid"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	path2 "path"

	"os"
	"time"
)

func UserRegister(context *gin.Context){
	DB:=common.GetDB()
	rdb,ctx:=common.GetRedis()
	var requestUser=model.User{}
	context.Bind(&requestUser)
	if requestUser.Email!="" {
		if err := gvalid.Check(requestUser.Email, "email", nil); err != nil {
			context.JSON(422,gin.H{"code":422,"msg":"该邮箱格式不对"})
			return
		}
	} else if requestUser.Phone!=""{
		if err := gvalid.Check(requestUser.Phone, "phone", nil); err != nil {
			context.JSON(422,gin.H{"code":422,"msg":"该手机号码格式不对"})
			return
		}
	}
	if err:=gvalid.Check(requestUser.Password,"password",nil);err!=nil{
		context.JSON(422,gin.H{"code":422,"msg":"密码格式不对"})
		return
	}
	if len(requestUser.NickName)==0{
		context.JSON(422,gin.H{"code":422,"msg":"昵称为空"})
		return
	} else if len(requestUser.NickName)>20{
		context.JSON(422,gin.H{"code":422,"msg":"昵称过长"})
		return
	}
	var user model.User
	//判断数据库是否有该账号
	DB.Where("email = ?", requestUser.Email).Or("phone = ?", requestUser.Phone).Find(&user)
	if user.Uuid!=""{
		context.JSON(422,gin.H{"code":422,"msg":"已有该账号"})
		return
	}

	hashPassword,err:=bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)

	if err!=nil{
		context.JSON(500,gin.H{"code":500,"msg":"系统密码处理错误"})
		return
	}
	u1 := uuid.NewV4()
	newUser:=model.User{
		Uuid: u1.String(),
		Email: requestUser.Email,
		Phone: requestUser.Phone,
		Password: string(hashPassword),

	}
	e:=DB.Create(&newUser)
	if e!=nil{
		if token, e:= common.ReleaseUserToken(newUser); e==nil{
			path:=fmt.Sprintf("public/user/%s",u1)
			os.MkdirAll(path, os.ModePerm)
			hset:=rdb.HSet(ctx,u1.String(),"token",token,"count",0,"time",time.Now().Unix())
			if hset.Err()!=nil{
				context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
				return
			}
			context.JSON(200,gin.H{"code":200,"msg":"注册成功","token":token})
		}else {
			//发放不了token
			context.JSON(500, gin.H{"code":500,"msg": "系统错误"})
		}
	}else {
		//添加不了进去数据库
		context.JSON(500, gin.H{"code":500,"msg": "系统错误"})
		panic(e)
	}
}
func UserLogin(context *gin.Context){
	DB:=common.GetDB()
	rdb,ctx:=common.GetRedis()
	var requestUser=model.User{}
	context.Bind(&requestUser)
	var user model.User
	DB.Where("email=?",requestUser.Email).Or("phone=?",requestUser.Phone).First(&user)


	if user.Uuid==""{
		context.JSON(422,gin.H{"code":422,"msg":"不存在该账号"})
		return
	}

	//判断是否冻结了
	t,err:=rdb.HGet(ctx,user.Uuid,"time").Int64()
	//fmt.Println(err)
	if err!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		panic(err)
		return
	}else if t>time.Now().Unix(){
		context.JSON(400, gin.H{"code": 400, "msg": "冻结","time":t})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password)) ; err != nil {
		context.JSON(400, gin.H{"code": 400, "msg": "密码错误"})
		//需要防止暴力破解

		count,err:=rdb.HGet(ctx,user.Uuid,"count").Int()
		if err!=nil{
			panic(err)
		}
		//fmt.Println(count)
		if count>=5{
			hset:=rdb.HSet(ctx,user.Uuid,"count",count+1,"time",time.Now().Unix()+int64((count-4)*60))
			panic(hset.Err())
		}else{
			hset:=rdb.HSet(ctx,user.Uuid,"count",count+1,"time",time.Now().Unix())
			panic(hset.Err())
		}
		return
	}

	token, err := common.ReleaseUserToken(user)
	if err != nil {
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})

		return
	}
	//将新的信息存入redis中

	hset:=rdb.HSet(ctx,user.Uuid,"token",token,"count",0,"time",time.Now().Unix())
	if hset.Err()!=nil{
		context.JSON(500, gin.H{"code": 500, "msg": "系统异常"})
		panic(hset.Err())
		return
	}
	//需要处理返回什么信息
	context.JSON(200,gin.H{"code":200,"msg":"登录成功","token":token,"data":model.ToUserDto(user)})
}
func UserUploadAvater(context *gin.Context){
	DB:=common.GetDB()
	user,_:=context.Get("user")
	u:=user.(model.User)
	avater,_:=context.FormFile("Avater")

	//判断是否图片，限制为jpg、bmp、jpeg、png
	//exts:=ExtAll(avater.Filename)
	//for i:=0;i<len(exts);i++{
	//	if exts[i]!=".jpg"&&exts[i]!=".jpeg"&&exts[i]!=".bmp"&&exts[i]!=".png"{
	//		context.JSON(400,gin.H{"code":400,"msg":"图片格式出错"})
	//		return
	//	}
	//}
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
	path:=fmt.Sprintf("public/user/%s/avater/",u.Uuid)
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
func UserUpdate(context *gin.Context){


}
func UserUpdatePassword(context *gin.Context){
	DB:=common.GetDB()
	user,_:=context.Get("user")
	u:=user.(model.User)
	password:=context.PostForm("Password")
	if err:=gvalid.Check(password,"password",nil);err!=nil{
		context.JSON(422,gin.H{"code":422,"msg":"密码格式不对"})
		return
	}
	hashPassword,_:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err:=DB.Model(&u).Update("password",hashPassword)
	//fmt.Println(err)
	//err:=DB.Exec("update user set password=? where uuid=?",hashPassword,u.Uuid)
	if err!=nil{
		context.JSON(200,gin.H{"code":200,"msg":"修改成功"})

	}else{
		context.JSON(500,gin.H{"code":500,"msg":"系统出错"})
	}


}
func UserForgetPassword(context *gin.Context){
	DB:=common.GetDB()
	user,_:=context.Get("user")
	u:=user.(model.User)
	password:=context.PostForm("Password")
	if err:=gvalid.Check(password,"password",nil);err!=nil{
		context.JSON(422,gin.H{"code":422,"msg":"密码格式不对"})
		return
	}
	hashPassword,_:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//fmt.Println(u.Uuid)
	err:=DB.Model(&u).Update("password",hashPassword)
	fmt.Println(err)
	if err!=nil{
		context.JSON(200,gin.H{"code":200,"msg":"修改成功"})

	}else{
		context.JSON(500,gin.H{"code":500,"msg":"系统出错"})
	}

}
func UserEmailVerify(context *gin.Context){
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
func UserPhoneVerify(context *gin.Context){

}
func UserInfo(context *gin.Context){
	user, _ := context.Get("user")
	//需要做返回信息处理
	context.JSON(200,gin.H{"code":200,"msg":"获取成功","data":model.ToUserDto(user.(model.User))})
}
//需要管理员权限
func UserDelete(context *gin.Context ){

}




//func ExtAll(path string) []string{
//	j:=len(path)
//	var ext []string
//	for i:=j-1;i>=0 && path[i]!= '/'; i--{
//		if path[i]=='.'{
//			ext=append(ext,path[i:j])
//			j=i
//		}
//	}
//	return ext
//}