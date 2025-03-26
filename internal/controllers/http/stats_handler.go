package http

import (
	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)


type StatHandler struct {
	service service.StatService
}

func NewStatHandler(service service.StatService) *StatHandler {
	return &StatHandler{
		service: service,
	}
}


func (sh *StatHandler) GetStats(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
	})
}