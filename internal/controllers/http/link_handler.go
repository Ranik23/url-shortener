package http

import (
	"context"
	"net/http"

	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)


type LinkHandler struct {
	service service.LinkService
}

func NewLinkHandler(service service.LinkService) *LinkHandler{
	return &LinkHandler{service: service}
}

func (lh *LinkHandler) CreateShortURL(c *gin.Context) {

	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	url, exists := req["url"]
	if !exists || url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	shortURL, err := lh.service.CreateShortURL(context.Background(), url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shortened_url": shortURL})
}


func (lh *LinkHandler) DeleteShortURL(c *gin.Context) {

	shortURL := c.Params.ByName("shorten_url")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}


	if err := lh.service.DeleteShortURL(context.Background(), shortURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}