package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kyomel/blog-management/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	categoryHandler *CategoryHandler,
	postHandler *PostHandler,
	tagHandler *TagHandler,
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

	posts := router.Group("/api/posts")
	{
		posts.GET("", postHandler.ListPosts)
		posts.GET("/:id", postHandler.GetPostByID)
		posts.GET("/slug/:slug", postHandler.GetPostBySlug)
	}

	tags := router.Group("/api/tags")
	{
		tags.GET("", tagHandler.ListTags)
		tags.GET("/:id", tagHandler.GetTagByID)
		tags.GET("/slug/:slug", tagHandler.GetTagBySlug)
		tags.GET("/:id/posts", tagHandler.GetPostsByTag)
	}

	posts.GET("/:id/tags", tagHandler.GetTagsByPost)

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

			adminPosts := admin.Group("/posts")
			{
				adminPosts.POST("", postHandler.CreatePost)
				adminPosts.PUT("/:id", postHandler.UpdatePost)
				adminPosts.DELETE("/:id", postHandler.DeletePost)
				adminPosts.PUT("/:id/publish", postHandler.PublishPost)
			}

			adminTags := admin.Group("/tags")
			{
				adminTags.POST("", tagHandler.CreateTag)
				adminTags.PUT("/:id", tagHandler.UpdateTag)
				adminTags.DELETE("/:id", tagHandler.DeleteTag)
			}

			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin dashboard"})
			})
		}
	}
}
