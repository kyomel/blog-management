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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
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
		&models.Comment{},
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
