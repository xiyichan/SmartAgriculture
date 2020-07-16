package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)
var ctx = context.Background()
var rdb *redis.Client


func InitRedis(){
	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"), // use default Addr
		Password: viper.GetString("redis.password"),               // no password set
		DB:       viper.GetInt("redis.db"),
		//可以加配置信息
		//DialTimeout:  10 * time.Second,
		//ReadTimeout:  30 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//PoolSize:     10,
		//PoolTimeout:  30 * time.Second,
	})
	_, err := rdb.Ping(ctx).Result()
	if err!=nil{
		panic("redis err:"+err.Error())
	}
}
func GetRedis() (*redis.Client,context.Context){
	return rdb,ctx
}