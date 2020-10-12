package controllers

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/encoding/gjson"
	"net/http"
	"server/common"
	"server/models"
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
	fmt.Print(response.GetHttpContentString())

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
	iotid := ctx.PostForm("IotId")
	order := fmt.Sprintf("{\"water_switch\":%v,\"fan_switch\":%v,\"light_switch\":%v}", waterSwitch, fanSwitch, lightSwitch)
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
	iotId := ctx.PostForm("IotId")
	userId := claims.(*common.Claims).ID
	var pi models.Pi
	err := db.Model(&pi).Where("iot_id=?", iotId).Update("user_id", userId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "绑定失败，数据库出错"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "绑定成功"})

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
