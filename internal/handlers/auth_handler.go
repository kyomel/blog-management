package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" || req.Fullname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, username, fullname, and password are required"})
		return
	}

	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "email already exists":
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		case "username already exists":
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case services.ErrUserNotActive:
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is not active"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), requestBody.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
