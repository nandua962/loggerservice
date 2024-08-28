package config

import (
	"fmt"
	"os"
	"partner/internal/entities"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig function
func LoadConfig(appName string) (*entities.EnvConfig, error) {

	var cfg entities.EnvConfig

	if _, err := os.Stat(".env"); err == nil {
		println("[ENV] Load env variables from .env")
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}

	}

	err := envconfig.Process(appName, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
