package http

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/Ranik23/url-shortener/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestCreateShortURL(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockLinkService := mock.NewMockLinkService(ctrl)


	mockLinkService.EXPECT().CreateShortURL(gomock.Any(), gomock.Any()).Return("", nil).MinTimes(1)

	linkHandler := NewLinkHandler(mockLinkService)

	router := gin.Default()
	router.POST("/api/shorten", linkHandler.CreateShortURL)


	req, _ := http.NewRequest("POST", "/api/shorten", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}	



// доделать 
func TestDeleteShortURL(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockLinkService := mock.NewMockLinkService(ctrl)


	mockLinkService.EXPECT().DeleteShortURL(gomock.Any(), gomock.Any()).Return("", nil).MinTimes(1)

	linkHandler := NewLinkHandler(mockLinkService)

	router := gin.Default()
	router.POST("/api/shorten", linkHandler.DeleteShortURL)


	req, _ := http.NewRequest("POST", "/api/shorten", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}	