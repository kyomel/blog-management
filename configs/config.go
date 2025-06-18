package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	// Enable replacement of '.' to '_' in env variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Check if .env file exists and load it
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Warning: Error reading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	// Set default values
	setDefaults()

	// Load YAML config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}

	// Map environment variables to configuration
	cfg := Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: viper.GetString("JWT_EXPIRY"),
		},
	}

	// Debug: Print configuration values (without sensitive data)
	log.Printf("Configuration values:")
	log.Printf("Server: port=%s, mode=%s", cfg.Server.Port, cfg.Server.Mode)
	log.Printf("Database: host=%s, port=%s, user=%s, dbname=%s, sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Check required configs
	if err := checkRequiredConfig(); err != nil {
		log.Printf("Configuration error: %v", err)
		return nil, err
	}

	return &cfg, nil
}

func checkRequiredConfig() error {
	required := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"JWT_SECRET",
	}

	for _, key := range required {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return fmt.Errorf("missing required configuration: %s (must be set via .env file or environment variable)", key)
		}
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")

	viper.SetDefault("JWT_EXPIRY", "24h")

}

func bindEnvVars() {
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")

	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.user", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.dbname", "DB_NAME")
	_ = viper.BindEnv("database.sslmode", "DB_SSLMODE")

	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")
	_ = viper.BindEnv("jwt.expiry", "JWT_EXPIRY")
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expiry string `mapstructure:"expiry"`
}
