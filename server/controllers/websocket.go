package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/common"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WSHandler(ctx *gin.Context) {
	//	w.Write([]byte("hello"))

	fmt.Println("connection websocket")
	var (
		wsConn *websocket.Conn
		err    error
		conn   *common.Connection
		data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		return
	}
	//TODO:传入用户id
	if conn, err = common.InitConnection(wsConn, "uuid"); err != nil {
		goto ERR
	}

	//TODO: 用户id作为key
	common.WebSocketMap.Range(func(k, v interface{}) bool {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
		return true
	})

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}
