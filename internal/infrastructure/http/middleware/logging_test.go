package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	e := echo.New()
	loggingMiddleware := NewLoggingMiddleware()
	e.Use(loggingMiddleware.Logger)

	e.POST("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "success"})
	})

	reqBody := map[string]string{"input": "test"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "test-value")

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	
	var respBody map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", respBody["message"])
}

func TestLoggingMiddleware_WithError(t *testing.T) {
	e := echo.New()
	loggingMiddleware := NewLoggingMiddleware()
	e.Use(loggingMiddleware.Logger)

	e.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}