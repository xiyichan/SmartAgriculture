package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"DEFAULT:NULL"`
	Phone    string `gorm:"unique;DEFAULT:NULL"`
	Email    string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avatar   string `gorm:"DEFAULT:NULL"`
}

func (a Admin) TableName() string {
	return "admin"
}

type AdminDto struct {
	Name   string `json:"Name"`
	Phone  string `json:"Phone"`
	Email  string `json:"Email"`
	Avatar string `json:"Avatar"`
}

func ToAdminDto(admin Admin) AdminDto {
	return AdminDto{
		Name:   admin.Name,
		Phone:  admin.Phone,
		Email:  admin.Email,
		Avatar: admin.Avatar,
	}
}
