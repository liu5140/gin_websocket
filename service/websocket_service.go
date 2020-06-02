package service

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	wsConn      *websocket.Conn
	inChan      chan []byte
	outChan     chan []byte
	mutex       sync.Mutex
	isClose     bool
	isCloseChan chan byte
}

var concurrentHolder sync.Map

func InitWebSocketClientService(userid int64, ws *WebSocketService) (err error) {
	if _, ok := concurrentHolder.Load(userid); !ok {
		concurrentHolder.Store(userid, ws)
	}
	return nil
}

//写数据到webScoket
func WriteMessage(userid int64, data []byte) (err error) {
	if syncma, ok := concurrentHolder.Load(userid); ok {
		ws := syncma.(*WebSocketService)
		ws.outChan <- data
	} else {
		//Log.Infoln("================没有建立连接=============")
	}
	return nil
}

func InitWebSocketService(c *gin.Context) (ws *WebSocketService, err error) {
	// change the reqest to websocket model
	wsconn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return nil, err
	}
	ws = &WebSocketService{
		wsConn:      wsconn,
		inChan:      make(chan []byte, 1000),
		outChan:     make(chan []byte, 1000),
		isCloseChan: make(chan byte, 1),
	}

	// if isNeedRead {
	// 	go webSocketService.read()
	// }
	// go webSocketService.writeText()
	return ws, nil
}

func (conn *WebSocketService) Read() (err error) {
	for {
		var data []byte
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			return err
		}
		//阻塞在这里，等到有人取消息
		select {
		case conn.inChan <- data:
		case <-conn.isCloseChan:
			conn.Close()
			return errors.New("connection is Close")
		}

	}
	return nil
}

//从outchan获取数据写回到websocket中
func (conn *WebSocketService) WriteText() (err error) {
	for {
		select {
		case message, ok := <-conn.outChan:
			//aa := string(message) + "返回"
			//Log.Infoln("===================", aa)
			if ok {
				if err = conn.wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
					//失败则关闭conn
					conn.Close()
					return err
				}
			}
		case <-conn.isCloseChan:
			conn.Close()
			return errors.New("connection is Close")
		}

	}
	return nil
}

//从outchan获取数据写回到websocket中
func (conn *WebSocketService) WriteTextToOutChann([]byte) (err error) {
	for {
		select {
		case message, ok := <-conn.outChan:
			if ok {
				if err = conn.wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
					//失败则关闭conn
					conn.Close()
					return err
				}
			}
		case <-conn.isCloseChan:
			conn.Close()
			return errors.New("connection is Close")
		}

	}
	return nil
}

//可多次调用
func (conn *WebSocketService) Close() (err error) {
	conn.wsConn.Close()
	//保证线程安全
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.isCloseChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
	return nil
}
