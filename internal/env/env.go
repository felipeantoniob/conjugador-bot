package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultEnvFilePath = ".env.local"

	botTokenKey = "BOT_TOKEN"
	guildIDKey  = "GUILD_ID"

	errLoadEnvFile = "failed to load environment file %s: %w"
	errMissingVars = "required environment variables are missing: %s"
)

// EnvLoader is an interface that abstracts the loading of environment variables.
type EnvLoader interface {
	Load(filePath string) error
}

// GodotenvLoader is an implementation of EnvLoader using the godotenv package.
type GodotenvLoader struct{}

func (l *GodotenvLoader) Load(filePath string) error {
	return godotenv.Load(filePath)
}

// LoadEnv loads environment variables from the default file path using the provided EnvLoader.
func LoadEnv(loader EnvLoader) error {
	return LoadEnvFromFile(defaultEnvFilePath, loader)
}

// LoadEnvFromFile loads environment variables from the specified file using the provided EnvLoader.
func LoadEnvFromFile(envFilePath string, loader EnvLoader) error {
	if err := loader.Load(envFilePath); err != nil {
		return fmt.Errorf(errLoadEnvFile, envFilePath, err)
	}
	return nil
}

// GetRequiredEnvVars retrieves required environment variables, returning an error if any are missing.
func GetRequiredEnvVars() (string, string, error) {
	botToken := os.Getenv(botTokenKey)
	guildID := os.Getenv(guildIDKey)

	var missingVars []string
	if botToken == "" {
		missingVars = append(missingVars, botTokenKey)
	}

	if guildID == "" {
		missingVars = append(missingVars, guildIDKey)
	}

	if len(missingVars) > 0 {
		return botToken, guildID, fmt.Errorf(errMissingVars, missingVars)
	}

	return botToken, guildID, nil
}
