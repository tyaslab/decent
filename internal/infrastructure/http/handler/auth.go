package handler

import (
	"decent/internal/application/dto"
	"decent/internal/infrastructure/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	jwtService *auth.JWTService
}

func NewAuthHandler(jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{jwtService: jwtService}
}

func (h *AuthHandler) GenerateToken(c echo.Context) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// if username == "" || password == "" {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	// }

	token, err := h.jwtService.GenerateToken(1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}
