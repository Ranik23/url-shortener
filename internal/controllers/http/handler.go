package http

import (
	"errors"
	"fmt"

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


func (h *Handler) AddRoute(method string, path string, fn gin.HandlerFunc) error {
	allowedMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
		"HEAD": true, "OPTIONS": true,
	}

	if method == "" || path == "" || fn == nil {
		return errors.New("invalid route parameters")
	}

	if !allowedMethods[method] {
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	h.Handle(method, path, fn)
	return nil
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