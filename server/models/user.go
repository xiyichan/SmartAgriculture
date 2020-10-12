package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"DEFAULT:NULL"`
	Phone    string `gorm:"unique;DEFAULT:NULL"`
	Email    string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avater   string `gorm:"DEFAULT:NULL"`

	Pis []Pi `gorm:"foreignKey:UserId"`
}
type UserDto struct {
	Name   string `json:"Name"`
	Phone  string `json:"hone"`
	Email  string `json:"Email"`
	Avater string `json:"Avater"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Name:   user.Name,
		Phone:  user.Phone,
		Email:  user.Email,
		Avater: user.Avater,
	}
}
