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

type GenerateTokenRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandler(jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{jwtService: jwtService}
}

func (h *AuthHandler) GenerateToken(c echo.Context) error {
	var request GenerateTokenRequestDto
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	}

	if request.Username != "admin" || request.Password != "password" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
	}

	token, err := h.jwtService.GenerateToken(1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}
