package main

import (
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json, text/javascript")
	chat := newChat(m)
	s.RegisterService(chat, "")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "client.html", gin.H{})
	})

	r.POST("/json", func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/channel/:name", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	r.Run(":8080")
}
