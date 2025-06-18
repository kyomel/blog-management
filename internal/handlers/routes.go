package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kyomel/blog-management/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	categoryHandler *CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	categories := router.Group("/api/categories")
	{
		categories.GET("", categoryHandler.ListCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	}

	api := router.Group("/api")
	api.Use(authMiddleware.Authenticate())
	{
		admin := api.Group("/admin")
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("", categoryHandler.CreateCategory)
				adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
				adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin dashboard"})
			})
		}
	}
}
