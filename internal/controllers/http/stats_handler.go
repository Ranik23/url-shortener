package http

import (
	"github.com/gin-gonic/gin"
)


type StatHandler struct {
	mainHandler *Handler
}

func NewStatHandler(mainHandler *Handler) *StatHandler {
	return &StatHandler{mainHandler: mainHandler}
}


func (sh *StatHandler) GetStats(c *gin.Context) {
	c.JSON(200, gin.H{"success": true})
}