package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadDotEnvDir() error {
	// Search for .env directory in the current working directory
	envDir := ".env"

	// Check if .env directory exists
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf(".env directory not found")
	}

	// Find all *.env files in the .env directory
	pattern := filepath.Join(envDir, "*.env")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to search for .env files: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no .env files found in %s directory", envDir)
	}

	// Load each .env file
	for _, envFile := range matches {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("failed to load %s: %w", envFile, err)
		}
	}

	return nil
}

func LoadConfigFromEnv() (*GlobalConfiguration, error) {
	config := &GlobalConfiguration{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	return config, nil
}
