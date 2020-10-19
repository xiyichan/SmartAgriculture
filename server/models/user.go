package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"DEFAULT:NULL"`
	Phone    string `gorm:"unique;DEFAULT:NULL"`
	Email    string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avatar   string `gorm:"DEFAULT:NULL"`

	Pis []Pi `gorm:"foreignKey:UserId"`
}
type UserDto struct {
	Id string`json:"Id"`
	Name   string `json:"Name"`
	Phone  string `json:"Phone"`
	Email  string `json:"Email"`
	Avatar string `json:"Avatar"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Name:   user.Name,
		Phone:  user.Phone,
		Email:  user.Email,
		Avatar: user.Avatar,
	}
}
