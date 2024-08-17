package envconfig

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

// LoadEnv loads environment variables from the default file.
func LoadEnv() error {
	return LoadEnvFromFile(defaultEnvFilePath)
}

// LoadEnvFromFile loads environment variables from the specified file.
func LoadEnvFromFile(envFilePath string) error {
	if err := godotenv.Load(envFilePath); err != nil {
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
		return "", "", fmt.Errorf(errMissingVars, missingVars)
	}

	return botToken, guildID, nil
}
