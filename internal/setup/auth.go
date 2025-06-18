package setup

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyomel/blog-management/internal/handlers"
	"github.com/kyomel/blog-management/internal/middleware"
	"github.com/kyomel/blog-management/internal/repositories"
	"github.com/kyomel/blog-management/internal/services"
	"github.com/kyomel/blog-management/internal/utils"
)

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func SetupAuth(router *gin.Engine, db *sql.DB, config AuthConfig) {
	userRepo := repositories.NewUserRepository(db)
	jwtService := utils.NewJWTService(
		config.AccessSecret,
		config.RefreshSecret,
		config.AccessExpiry,
		config.RefreshExpiry,
	)

	authService := services.NewAuthService(
		userRepo,
		jwtService,
		config.AccessExpiry,
	)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	authHandler := handlers.NewAuthHandler(authService)

	handlers.RegisterRoutes(router, authHandler)

	protected := router.Group("/api")
	protected.Use(authMiddleware.Authenticate())
	{
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin dashboard"})
			})
		}
	}
}
