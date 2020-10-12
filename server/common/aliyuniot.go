package common

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/spf13/viper"
)

var Client *sdk.Client

func InitAliyunIot() {
	regionId := viper.GetString("aliyuniot.regionId")
	accessKey := viper.GetString("aliyuniot.accessKey")
	accessSecret := viper.GetString("aliyuniot.accessSecret")

	client, err := sdk.NewClientWithAccessKey(regionId, accessKey, accessSecret)
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	Client = client
}
func GetAliyunIotClient() *sdk.Client {
	return Client
}
