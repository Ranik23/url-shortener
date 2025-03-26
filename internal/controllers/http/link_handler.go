package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)


type LinkHandler struct {
	service service.LinkService
}

func NewLinkHandler(service service.LinkService) *LinkHandler{
	return &LinkHandler{
		service: service,
	}
}

func (lh *LinkHandler) CreateShortURL(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
	})
}


func (lh *LinkHandler) DeleteShortURL(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
	})
}