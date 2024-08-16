package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	botTokenKey = "BOT_TOKEN"
	guildIDKey  = "GUILD_ID"

	errMsgLoadEnvFile = "failed to load environment file %s: %w"
	errMsgMissingVars = "required environment variables are missing: %s"
)

var defaultEnvFilePath = ".env.local"

// LoadEnv loads environment variables from the specified file.
func LoadEnv(envFilePath ...string) error {
	filePath := defaultEnvFilePath
	if len(envFilePath) > 0 {
		filePath = envFilePath[0]
	}

	if err := godotenv.Load(filePath); err != nil {
		return fmt.Errorf(errMsgLoadEnvFile, filePath, err)
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
		return "", "", fmt.Errorf(errMsgMissingVars, missingVars)
	}

	return botToken, guildID, nil
}
