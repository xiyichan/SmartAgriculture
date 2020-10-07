package common

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"time"
)

var logger *logrus.Logger
var logBuf *LogBuffer
var logCollection *mongo.Collection

type LogBuffer struct {
	io.Writer
	logChan chan []byte
}

func (w *LogBuffer) Write(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, errors.New("no data")
	}
	w.logChan <- p
	return len(p), nil
}


func InitLog() {
	logger = logrus.New()
	logBuf = &LogBuffer{
		logChan: make(chan []byte, 100),
	}
	logger.SetOutput(logBuf)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logCollection = GetMango().Database(viper.GetString("mongodb.logdb")).Collection(viper.GetString("mongodb.logcollection"))
	go StartLogSystem()
}

func StartLogSystem() {
	var logDataDocs []interface{}
	for {
		select {
		case logByte := <-logBuf.logChan:
			logDataDoc := &bson.M{}
			err := json.Unmarshal(logByte, logDataDoc)
			if err != nil {
				panic(err)
			}else {
				logDataDocs=append(logDataDocs,logDataDoc)
				if len(logDataDocs)>=10{
					ctx, cancel := context.WithTimeout(context.Background(),time.Second * 10)
					_, err2 := logCollection.InsertMany(ctx, logDataDocs, options.InsertMany())
					if err2 != nil {
						cancel()
						panic(err2)
					}
					logDataDocs = make([]interface{}, 0,11)
				}
			}
		}
	}
}


func LogErr(msg string, code string) {
	//TODO:写入上下文用户id
	//login,_:=context.Get("loginer")
	//loginer:=login.(model.jwtLoginer)
	logger.WithFields(logrus.Fields{
		"status_code":code,
		"req_uri":"context.Request.URL",
		"req_method":"context.Request.Method",
		"client_ip":"context.ClientIP()",
		"user":"loginer.Uuid",
	}).Error(msg)
}



func LogWarning(msg string, code string) {
	//TODO:写入上下文用户id
	//login,_:=context.Get("loginer")
	//loginer:=login.(model.jwtLoginer)
	logger.WithFields(logrus.Fields{
		"status_code": code,
		"req_uri":     "context.Request.URL",
		"req_method":  "context.Request.Method",
		"client_ip":   "context.ClientIP()",
		"user":        "loginer.Uuid",
	}).Warn(msg)
}
