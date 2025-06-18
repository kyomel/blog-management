package database

import (
	"fmt"
	"log"

	"github.com/kyomel/blog-management/configs"
	"github.com/kyomel/blog-management/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(config *configs.DatabaseConfig) error {
	if config.Host == "" || config.User == "" || config.Password == "" || config.DBName == "" {
		return fmt.Errorf("missing required database configuration (host, user, password, or dbname)")
	}

	port := config.Port
	if port == "" {
		port = "5432"
	}

	sslmode := config.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		port,
		sslmode,
	)

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		// Return a more descriptive error
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	log.Println("Starting database migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.MediaFile{},
		&models.AuditLog{},
	)

	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
