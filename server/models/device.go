package models

import (
	"gorm.io/gorm"
)

type Pi struct {
	gorm.Model
	ProductKey   string
	DeviceName   string `gorm:"unique;not null;primary_key"`
	DeviceSecret string
	NickName     string
	IotId        string `gorm:"unique;not null;primary_key"`

	Status      string
	FanSwitch   *bool
	WaterSwitch *bool
	LightSwitch *bool

	Temperature    float32
	Humidity       float32
	SoilMoisture   int
	LightIntensity int

	UserId uint `gorm:"foreignKey:UserId"`
}

//返回当前数据给用户的
type PiData struct {
	IotId          string
	Status         string
	FanSwitch      bool
	WaterSwitch    bool
	LightSwitch    bool
	Temperature    float32
	Humidity       float32
	SoilMoisture   int
	LightIntensity int
}

type PiHistoryData struct {
	gorm.Model
	IotId          string
	Temperature    float32
	Humidity       float32
	SoilMoisture   int
	LightIntensity int
	Time           string
}
