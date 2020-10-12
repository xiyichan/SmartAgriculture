package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"DEFAULT:NULL"`
	Phone    string `gorm:"unique;DEFAULT:NULL"`
	Email    string `gorm:"unique;DEFAULT:NULL"`
	Password string `gorm:"DEFAULT:NULL"`
	Avater   string `gorm:"DEFAULT:NULL"`
}
