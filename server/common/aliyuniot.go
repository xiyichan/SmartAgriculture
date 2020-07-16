package common

import (

"github.com/aliyun/alibaba-cloud-sdk-go/sdk"

"github.com/spf13/viper"
)
var Client *sdk.Client

func InitIot() *sdk.Client{
	regionId:=viper.GetString("iotsource.regionId")
	accessKey:=viper.GetString("iotsource.accessKey")
	accessSecret:=viper.GetString("iotsource.accessSecret")

	client,err:=sdk.NewClientWithAccessKey(regionId,accessKey,accessSecret)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	Client=client

	return Client
}
func GetClient() *sdk.Client{
	return Client
}