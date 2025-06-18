package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes sets up all the auth routes
func RegisterRoutes(router *gin.Engine, authHandler *AuthHandler) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}
}
