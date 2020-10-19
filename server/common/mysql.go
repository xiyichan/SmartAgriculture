package common

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"server/models"
	"strconv"
)

var DB *gorm.DB

func InitMysql() {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	//自动生成数据库表
	db.AutoMigrate(&models.Pi{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Admin{})

	DB = db
	//初始化admin账号
	var admin models.Admin
	db.Find(&admin).Where("email=?", "2863768433@qq.com")
	if admin.ID == 0 {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte("czx987852"), bcrypt.DefaultCost)
		newAdmin := models.Admin{
			Email:    "2863768433@qq.com",
			Password: string(hashPassword),
			Name:     "chenlifeng",
		}
		e := db.Create(&newAdmin).Error
		if e != nil {
			panic(e)
		}
	}

	path := fmt.Sprintf("public/admin/%d/", admin.ID)
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	//year:=time.Now().Year()
	//month:=time.Now().Month()
	for i := 2020; i < 2022; i++ {
		for j := 1; j <= 12; j++ {
			piname := strconv.Itoa(i) + "-" + strconv.Itoa(j) + "-pi"
			if db.Migrator().HasTable(piname) == false {
				db.Migrator().CreateTable(&models.PiHistoryData{})
				db.Migrator().RenameTable(&models.PiHistoryData{}, piname)
			}

		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
