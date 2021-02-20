package common

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var mgoCli *mongo.Client
var mgoFirst sync.Once

func GetMongoClient() (*mongo.Client, error) {
	var err error
	mgoFirst.Do(func() {
		host := viper.GetString("mongodb.host")
		port := viper.GetString("mongodb.port")
		database := viper.GetString("mongodb.database")
		username := viper.GetString("mongodb.username")
		password := viper.GetString("mongodb.password")
		//database要选admin才行
		args := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?w=majority", username, password, host, port, database)
		clientOptions := options.Client().ApplyURI(args)
		mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	})
	if err != nil {

		return nil, err
	}
	return mgoCli, nil
}
