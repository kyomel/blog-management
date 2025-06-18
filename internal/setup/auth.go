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
	categoryRepo := repositories.NewCategoryRepository(db)

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

	categoryService := services.NewCategoryService(categoryRepo)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Register all routes with appropriate middleware
	handlers.RegisterRoutes(router, authHandler, categoryHandler, authMiddleware)
}
