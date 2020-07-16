package controller

import (
	"server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"fmt"
)

func DeviceSetPower(ctx *gin.Context){
	client:=common.GetClient()
	powerSwitch:=ctx.PostForm("PowerSwitch")
	iotid:=ctx.PostForm("IotId")
	order:=fmt.Sprintf("{\"PowerSwitch\":%v}",powerSwitch)
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
	ctx.JSON(http.StatusOK,gin.H{"code":200,"data":response.GetHttpContentString(),"msg":"控制成功"})
}