package common

import (
	"server/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	//命名按照原本得
	db.SingularTable(true)
	//自动生成数据库表


	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Admin{})




	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}