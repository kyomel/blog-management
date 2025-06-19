package main

import (
	"log"
	"time"

	"github.com/kyomel/blog-management/configs"
	"github.com/kyomel/blog-management/internal/database"
	"github.com/kyomel/blog-management/internal/setup"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := database.Connect(&config.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	db, err := database.GetDB().DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	gin.SetMode(config.Server.Mode)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog Management API is running",
		})
	})

	accessExpiry, err := time.ParseDuration(config.JWT.AccessExpiry)
	if err != nil {
		log.Printf("Warning: Invalid JWT access expiry format, using default 15m: %v", err)
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(config.JWT.RefreshExpiry)
	if err != nil {
		log.Printf("Warning: Invalid JWT refresh expiry format, using default 7d: %v", err)
		refreshExpiry = 7 * 24 * time.Hour
	}

	setup.SetupAuth(router, db, setup.AuthConfig{
		AccessSecret:  config.JWT.AccessSecret,
		RefreshSecret: config.JWT.RefreshSecret,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
		Cloudinary:    config.Cloudinary,
	})

	log.Printf("Server starting on port %s", config.Server.Port)
	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
