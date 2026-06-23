package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/balitravelhealth/platform/gateway-go/internal/services"
)

type adminLoginRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	DeviceInfo string `json:"device_info"`
}

type bootstrapRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// POST /admin/auth/login — GO-23
func (h *Handler) AdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	pair, user, err := h.AdminService.AdminLogin(c.Request.Context(), req.Email, req.Password, req.DeviceInfo, h.AuthService)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// POST /admin/bootstrap — only works when no admin exists yet
func (h *Handler) AdminBootstrap(c *gin.Context) {
	var req bootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password (min 8 chars) are required"})
		return
	}

	pair, user, err := h.AdminService.Bootstrap(c.Request.Context(), req.Email, req.Password, h.AuthService)
	if err != nil {
		if err.Error() == "admin sudah ada" {
			c.JSON(http.StatusConflict, gin.H{"error": "admin already exists"})
			return
		}
		if errors.Is(err, services.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
