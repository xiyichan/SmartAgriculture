package common

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mgoCli *mongo.Client

func InitMongodb() {
	var err error
	host := viper.GetString("mongodb.host")
	port := viper.GetString("mongodb.port")
	database := viper.GetString("mongodb.logdb")
	username := viper.GetString("mongodb.username")
	password := viper.GetString("mongodb.password")
	args := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", username, password, host, port, database)
	clientOptions := options.Client().ApplyURI(args)

	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}

func GetMango() *mongo.Client {
	return mgoCli
}
func GetCollection(collection string) *mongo.Collection {
	c := mgoCli.Database(viper.GetString("mongodb.logdb")).Collection(collection)
	return c
}
