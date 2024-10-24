package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		i := 0
		for {
			i++
			conn.WriteMessage(websocket.TextMessage, []byte("New message (#"+strconv.Itoa(i)+")"))
			time.Sleep(time.Second)
		}
	})
	router.Run(":8002")
}
