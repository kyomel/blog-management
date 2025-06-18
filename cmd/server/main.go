package main

import (
	"log"

	"github.com/kyomel/blog-management/configs"
	"github.com/kyomel/blog-management/internal/database"

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

	gin.SetMode(config.Server.Mode)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog Management API is running",
		})
	})

	log.Printf("Server starting on port %s", config.Server.Port)
	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
