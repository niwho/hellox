package http

import (
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
