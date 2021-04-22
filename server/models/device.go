package models

import (
	"gorm.io/gorm"
)

type Pi struct {
	gorm.Model
	ProductKey     string
	DeviceName     string `gorm:"unique;not null;primary_key"`
	DeviceSecret   string
	NickName       string
	IotId          string `gorm:"unique;not null;primary_key"`
	Status         string
	FanSwitch      *bool
	WaterSwitch    *bool
	LightSwitch    *bool
	AutoSwitch     *bool
	Temperature    float32
	Humidity       float32
	SoilMoisture   int
	LightIntensity int
	UserId         uint `gorm:"foreignKey:UserId"`
}

func (a Pi) TableName() string {
	return "pi"
}

type PiListData struct {
	DeviceName   string
	DeviceSecret string
	NickName     string
	IotId        string
	Status       string
	UserId       uint
}

//返回当前数据给用户的
type PiData struct {
	IotId          string
	Status         string
	FanSwitch      bool
	WaterSwitch    bool
	LightSwitch    bool
	AutoSwitch     bool
	Temperature    float32
	Humidity       float32
	SoilMoisture   int
	LightIntensity int
}

type PiHistoryData struct {
	IotId          string  `bson:"iot_id"`
	Temperature    float32 `bson:"temperature"`
	Humidity       float32 `bson:"humidity"`
	SoilMoisture   int     `bson:"soil_moisture"`
	LightIntensity int     `bson:"light_intensity"`
	Time           int64   `bson:"time"`
}
