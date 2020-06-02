package router

import (
	"gin_websocket/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func webSocketRouter(r *gin.RouterGroup) {
	r.GET("/ws", getWebSocket)
	r.GET("/ws/response", responseWebSocket)
}

func responseWebSocket(c *gin.Context) {
	//	profile := c.MustGet(PROFILE).(middleware.Profile)
	message, _ := c.GetQuery("message")
	//	Log.Infoln("==============", message)
	err := service.WriteMessage(11, []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "错误")
		return
	}
	c.JSON(http.StatusOK, "错误")
}

func getWebSocket(c *gin.Context) {
	//解析token 获取用户id
	//token, ok := c.GetQuery("token")
	//Log.Infoln("=================", token)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, "token错误")
	// 	return
	// }

	// session := service.NewSessionService()
	// id, err := session.ParseToken(token, hmacSampleSecret)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "错误")
	// 	return
	// }
	var id int64
	id = 1
	wsconn, err := service.InitWebSocketService(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "错误")
		return
	}

	err = service.InitWebSocketClientService(id, wsconn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "错误")
		return
	}
	go wsconn.Read()
	go wsconn.WriteText()
}
