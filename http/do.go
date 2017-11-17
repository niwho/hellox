package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/niwho/hellox/im"
)

type Resp struct {
	Status  int
	Message string
}

func stat(c *gin.Context) {
	c.JSON(200, Resp{Status: 0, Message: "tbd"})
}

func serveWS(c *gin.Context) {
	im.ServeWs(c.Writer, c.Request)
}

func broadcast(c *gin.Context) {
	text := c.Query("text")
	text = strings.TrimSpace(text)
	if text != "" {
		im.Broadcast(text)
	}
}
