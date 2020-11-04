package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/encoding/gjson"
	"net/http"
	"server/common"
	"server/models"
	"strconv"
)

func RegisterPi(ctx *gin.Context) {
	claims, _ := ctx.Get("Claims")
	//fmt.Println(claims)
	if claims.(*common.Claims).Role != "admin" {
		ctx.JSON(422, gin.H{"code": 422, "msg": "权限不足"})
		return
	}
	client := common.GetAliyunIotClient()
	db := common.GetDB()
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "iot.cn-shanghai.aliyuncs.com"
	request.Version = "2018-01-20"
	request.ApiName = "RegisterDevice"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["ProductKey"] = "a1I7rPDpEx5"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	//fmt.Print(response.GetHttpContentString())

	j := gjson.New(response.GetHttpContentString())
	data := gjson.New(j.Get("Data"))
	IotId := data.GetString("IotId")
	DeviceSecret := data.GetString("DeviceSecret")
	ProductKey := data.GetString("ProductKey")
	DeviceName := data.GetString("DeviceName")

	newPi := models.Pi{
		IotId:        IotId,
		DeviceSecret: DeviceSecret,
		ProductKey:   ProductKey,
		DeviceName:   DeviceName,
		NickName:     "pi",
	}

	//qr,err:=qrcode.New(IotId,qrcode.Medium)
	//if err != nil {
	//	//log.Fatal(err)
	//	panic(err)
	//} else {
	//	qr.BackgroundColor = color.RGBA{255,255,255,255}
	//	qr.ForegroundColor = color.White
	//	qr.WriteFile(256,"./golang_qrcode.png")
	//}
	e := db.Create(&newPi).Error
	if e != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "系统错误，数据库异常"})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "msg": "新建设备成功"})

}
func UserSetPiProperty(ctx *gin.Context) {
	client := common.GetAliyunIotClient()
	//powerSwitch:=ctx.PostForm("PowerSwitch")

	waterSwitch := ctx.PostForm("WaterSwitch")
	fanSwitch := ctx.PostForm("FanSwitch")
	lightSwitch := ctx.PostForm("lightSwitch")
	var w,f,l int
	if waterSwitch=="true"{
		w=1
	}
	if fanSwitch=="true"{
		f=1
	}
	if lightSwitch=="true"{
		l=1
	}
	iotid := ctx.PostForm("IotId")
	//fmt.Println(waterSwitch)
	order := fmt.Sprintf("{\"water_switch\":%v,\"fan_switch\":%v,\"light_switch\":%v}", w, f, l)
	//order:="{\"water_switch\":1}"
	//fmt.Println(order)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "iot.cn-shanghai.aliyuncs.com"
	request.Version = "2018-01-20"
	request.ApiName = "SetDeviceProperty"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["Items"] = order
	request.QueryParams["IotId"] = iotid
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	//fmt.Print(response.GetHttpContentString())
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": response.GetHttpContentString(), "msg": "控制成功"})
}
func UserBindPi(ctx *gin.Context) {
	db := common.GetDB()
	claims, _ := ctx.Get("Claims")
	//iotId := ctx.PostForm("IotId")
	deviceName:=ctx.PostForm("DeviceName")
	userId := claims.(*common.Claims).ID
	//var pi models.Pi
	client := common.GetAliyunIotClient()

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "iot.cn-shanghai.aliyuncs.com"
	request.Version = "2018-01-20"
	request.ApiName = "RRpc"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DeviceName"] = "d5R46fOSfNNwVTNRuSaM"
	request.QueryParams["Timeout"] = "5000"

	msg:=base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(userId))))
	request.QueryParams["RequestBase64Byte"] = msg
	request.QueryParams["ProductKey"] = "a1I7rPDpEx5"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	//fmt.Print(response.GetHttpContentString())
	j := gjson.New(response.GetHttpContentString())
	rrpcCode:=j.GetString("RrpcCode")
	fmt.Println(rrpcCode)
	if rrpcCode=="SUCCESS" {
		fmt.Println(deviceName)
		e1:= db.Model(&models.Pi{}).Where("device_name=?", deviceName).Update("user_id", userId).Error
		if e1!= nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "绑定失败，数据库出错"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "绑定成功"})

	}else{
		ctx.JSON(500, gin.H{"code": 500, "msg": "设备有问题"})
	}
	//TODO:tips需要修改

}
func UserGetPiProperty(ctx *gin.Context) {
	db := common.GetDB()
	iotId := ctx.PostForm("IotId")
	var pi models.Pi
	var piData models.PiData
	pi.IotId = iotId
	err := db.Model(&pi).Find(&piData).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库出错"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "查询成功", "data": piData})
}

func UserGetPiHistoryData(ctx *gin.Context) {
	//claims,_:=ctx.Get("Claims")
	//ctx.PostForm("IotId")
	//userId:=claims.(*common.Claims).ID

}
func PiList(ctx *gin.Context) {
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
	var dto []models.PiListData

	var count int64
	db.Model(&models.Pi{}).Count(&count)
	e := db.Model(&models.Pi{}).Offset((pi - 1) * ps).Limit(ps).Find(&dto).Error
	if e != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "查询失败"})
	} else {
		ctx.JSON(200, gin.H{"code": 200, "data": dto, "msg": "查询成功", "count": count})
	}
}
func UserPiList(ctx *gin .Context){
	claims, _ := ctx.Get("Claims")
	c := claims.(*common.Claims)
	//fmt.Println(c)
	//fmt.Println(c.ID)
	userId:=c.ID
	db := common.GetDB()
	pageIndex := ctx.PostForm("PageIndex")
	pageSize := ctx.PostForm("PageSize")
	if pageIndex == "" || pageSize == "" {
		ctx.JSON(422, gin.H{"code": 422, "msg": "分页设置为空"})
		return
	}
	pi, _ := strconv.Atoi(pageIndex)
	ps, _ := strconv.Atoi(pageSize)
	var dto []models.PiData

	var count int64
	//fmt.Println(userId)
	db.Model(&models.Pi{}).Where("user_id=?",userId).Count(&count)
	e := db.Model(&models.Pi{}).Offset((pi - 1) * ps).Limit(ps).Where("user_id=?",userId).Find(&dto).Error
	if e != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "查询失败"})
	} else {
		ctx.JSON(200, gin.H{"code": 200, "data": dto, "msg": "查询成功", "count": count})
	}
}

func PiDelete(ctx *gin.Context) {
	claims, _ := ctx.Get("Claims")
	//fmt.Println(claims.(common.Claims).Role)
	if claims.(*common.Claims).Role != "admin" {
		ctx.JSON(422, gin.H{"code": 200, "msg": "权限不足"})
		return
	}

	client := common.GetAliyunIotClient()
	iotId := ctx.PostForm("IotId")
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "iot.cn-shanghai.aliyuncs.com"
	request.Version = "2018-01-20"
	request.ApiName = "DeleteDevice"
	request.QueryParams["RegionId"] = "cn-hangzhou"

	request.QueryParams["IotId"] = iotId

	_, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	//fmt.Println(response)

	db := common.GetDB()

	//var pi models.Pi
	//pi.IotId=iotId
	e := db.Where("iot_id=?", iotId).Delete(&models.Pi{}).Error
	if e != nil {
		ctx.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
	} else {
		ctx.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
	}

}
//设备登出
func PiLoginOut(ctx *gin.Context){
	db := common.GetDB()
	//iotId := ctx.PostForm("IotId")
	deviceName:=ctx.PostForm("DeviceName")
	deviceSecret:=ctx.PostForm("DeviceSecret")
	fmt.Println(deviceSecret,deviceSecret)

	e1:= db.Model(&models.Pi{}).Where("device_name=?", deviceName).Where("device_secret=?",deviceSecret).Update("user_id", 0).Error
	if e1!= nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "解绑失败，数据库出错"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "解绑成功"})
}
//用户客户端登出，还需要rrpc

func PiUserLoginOut(ctx *gin.Context){

}