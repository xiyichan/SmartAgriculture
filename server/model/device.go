package model

type Device struct {
	IotId string `gorm:"unique;not null;primary_key;auto_increment:false"`
	ProductKey string `gorm:"not null"`
	DeviceName string	`gorm:"not null"`
	DeviceSecret string	`gorm:"not null"`
	NickName string
	Status int
	PowerSwitch int
	Data string


}

