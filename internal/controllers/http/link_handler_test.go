package http

// import (
// 	"bytes"
// 	"encoding/json"
// 	serviceMock "github.com/Ranik23/url-shortener/internal/service/mock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestCreateShortURL(t *testing.T) {
// 	mockLinkService := serviceMock.NewLinkService(t)

// 	// Устанавливаем ожидание вызова метода CreateShortURL с конкретным URL
// 	expectedURL := "https://example.com"
// 	mockLinkService.
// 		On("CreateShortURL", mock.Anything, expectedURL).
// 		Return("shortened", nil)

// 	linkHandler := NewLinkHandler(mockLinkService)

// 	router := gin.Default()
// 	router.POST("/api/shorten", linkHandler.CreateShortURL)

// 	body := map[string]string{"url": expectedURL}
// 	jsonBody, err := json.Marshal(body)
// 	assert.NoError(t, err)

// 	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)

// 	mockLinkService.AssertExpectations(t)
// }

// func TestDeleteShortURL(t *testing.T) {
// 	mockLinkService := serviceMock.NewLinkService(t)

// 	mockLinkService.
// 		On("DeleteShortURL", mock.Anything, "abc123").
// 		Return(nil)

// 	linkHandler := NewLinkHandler(mockLinkService)

// 	router := gin.Default()
// 	router.POST("/api/delete/:shorten_url", linkHandler.DeleteShortURL)

// 	req, _ := http.NewRequest("POST", "/api/delete/abc123", nil)
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)

// 	mockLinkService.AssertExpectations(t)
// }
