package common

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

var WebSocketMap sync.Map

type Connection struct{
	uuid string
	wsConnect *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex  // 对closeChan关闭上锁
	isClosed bool  // 防止closeChan被关闭多次
}

func InitConnection(wsConn *websocket.Conn,uuid string)(conn *Connection ,err error){
	conn = &Connection{
		uuid: uuid,
		wsConnect:wsConn,
		inChan: make(chan []byte,10),
		outChan: make(chan []byte,10),
		closeChan: make(chan byte,1),
	}
	// 启动读协程
	go conn.readLoop();
	// 启动写协程
	go conn.writeLoop();
	return
}

func (conn *Connection)ReadMessage()(data []byte , err error){
	select{
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection)WriteMessage(data []byte)(err error){
	select{
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection)Close(){
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
		fmt.Println("close")
		WebSocketMap.Delete(conn.uuid)
		WebSocketMap.Range(func(k,v interface{})bool{
			fmt.Print(k)
			fmt.Print(":")
			fmt.Print(v)
			fmt.Println()
			return true
		})
	}
	conn.mutex.Unlock()
}

// 内部实现
func (conn *Connection)readLoop(){
	var(
		data []byte
		err error
	)
	for{
		//从连接中读数据
		if _, data , err = conn.wsConnect.ReadMessage(); err != nil{
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select{
		case conn.inChan <- data:
		case <- conn.closeChan:		// closeChan 感知 conn断开
			goto ERR
		}

	}
ERR:
	conn.Close()
}

func (conn *Connection)writeLoop(){
	var(
		data []byte
		err error
	)
	for{
		select{
		case data= <- conn.outChan:
			if err = conn.wsConnect.WriteMessage(websocket.TextMessage , data); err != nil{
				goto ERR
			}
		case <- conn.closeChan:
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection)checkConnection(){
	pingHandler := conn.wsConnect.PingHandler()
	err:=pingHandler("check connection")
	if err!=nil{
		conn.Close()
	}
}