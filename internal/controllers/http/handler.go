package http

import (
	"errors"
	"fmt"
	"log"

	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)


type Handler struct {
	*gin.Engine
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		Engine: gin.Default(),
		service: service,
	}
}


func (h *Handler) AddRoute(method, path string, fn gin.HandlerFunc) error {
	if method == "" || path == "" || fn == nil {
		return errors.New("invalid route parameters")
	}

	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic in AddRoute: %v\n", err)
			}
		}()

		h.Handle(method, path, fn)
		return nil

	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}
}


func (h *Handler) SetUpRoutes() {

	linkHandler := NewLinkHandler(h)
	statHandler := NewStatHandler(h)

	api := h.Group("/api") 
	{
		api.POST("/shorten", linkHandler.CreateShortURL)
		api.DELETE("/delete/:shorten_url", linkHandler.DeleteShortURL)
		api.GET("/stats/:shorten_url", statHandler.GetStats)
	}
}