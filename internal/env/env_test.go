package env

import (
	"fmt"
	"os"
	"testing"
)

const (
	botTokenValue = "test-bot-token"
	guildIDValue  = "test-guild-id"
)

// MockEnvLoader is a simple mock implementation of the EnvLoader interface for testing purposes.
type MockEnvLoader struct {
	LoadFunc func(filePath string) error
}

func (m *MockEnvLoader) Load(filePath string) error {
	if m.LoadFunc != nil {
		return m.LoadFunc(filePath)
	}
	return nil
}

// TestLoadEnv tests the LoadEnv function.
func TestLoadEnv(t *testing.T) {
	mockLoader := &MockEnvLoader{
		LoadFunc: func(filePath string) error {
			if filePath == defaultEnvFilePath {
				return nil
			}
			return fmt.Errorf("unexpected file path: %s", filePath)
		},
	}

	err := LoadEnv(mockLoader)
	if err != nil {
		t.Fatalf("LoadEnv() returned an error: %v", err)
	}
}

// TestLoadEnv_ErrorHandling checks error handling in LoadEnv.
func TestLoadEnv_ErrorHandling(t *testing.T) {
	mockLoader := &MockEnvLoader{
		LoadFunc: func(filePath string) error {
			return fmt.Errorf("some error")
		},
	}

	err := LoadEnv(mockLoader)
	if err == nil {
		t.Fatal("LoadEnv() expected an error, got nil")
	}
	if !contains(err.Error(), "failed to load environment file .env.local") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestLoadEnvFromFile tests the LoadEnvFromFile function.
func TestLoadEnvFromFile(t *testing.T) {
	mockLoader := &MockEnvLoader{
		LoadFunc: func(filePath string) error {
			if filePath == ".env.local" {
				return nil
			}
			return fmt.Errorf("unexpected file path: %s", filePath)
		},
	}

	err := LoadEnvFromFile(".env.local", mockLoader)
	if err != nil {
		t.Fatalf("LoadEnvFromFile() returned an error: %v", err)
	}
}

// TestGodotenvLoader tests the GodotenvLoader implementation.
func TestGodotenvLoader(t *testing.T) {
	loader := &GodotenvLoader{}
	envFile := ".env.test"
	defer func() {
		if err := os.Remove(envFile); err != nil {
			t.Fatalf("Failed to remove test file: %v", err)
		}
	}()

	if err := os.WriteFile(envFile, []byte("TEST_VAR=value"), 0666); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	if err := loader.Load(envFile); err != nil {
		t.Fatalf("GodotenvLoader.Load() returned an error: %v", err)
	}

	if got := os.Getenv("TEST_VAR"); got != "value" {
		t.Errorf("os.Getenv(\"TEST_VAR\") = %q; want %q", got, "value")
	}
}

// clearEnvVars clears specific environment variables for testing.
func clearEnvVars() {
	os.Unsetenv(botTokenKey)
	os.Unsetenv(guildIDKey)
}

func TestGetRequiredEnvVars(t *testing.T) {
	tests := []struct {
		name             string
		setEnvVars       func()
		expectedBotToken string
		expectedGuildID  string
		expectedError    string
	}{
		{
			name: "All environment variables set",
			setEnvVars: func() {
				os.Setenv(botTokenKey, "test-bot-token")
				os.Setenv(guildIDKey, "test-guild-id")
			},
			expectedBotToken: "test-bot-token",
			expectedGuildID:  "test-guild-id",
			expectedError:    "",
		},
		{
			name: "Bot token missing",
			setEnvVars: func() {
				os.Unsetenv(botTokenKey)
				os.Setenv(guildIDKey, "test-guild-id")
			},
			expectedBotToken: "",
			expectedGuildID:  "test-guild-id",
			expectedError:    "required environment variables are missing: [BOT_TOKEN]",
		},
		{
			name: "Guild ID missing",
			setEnvVars: func() {
				os.Setenv(botTokenKey, "test-bot-token")
				os.Unsetenv(guildIDKey)
			},
			expectedBotToken: "test-bot-token",
			expectedGuildID:  "",
			expectedError:    "required environment variables are missing: [GUILD_ID]",
		},
		{
			name: "Both environment variables missing",
			setEnvVars: func() {
				os.Unsetenv(botTokenKey)
				os.Unsetenv(guildIDKey)
			},
			expectedBotToken: "",
			expectedGuildID:  "",
			expectedError:    "required environment variables are missing: [BOT_TOKEN GUILD_ID]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnvVars()
			defer clearEnvVars()

			botToken, guildID, err := GetRequiredEnvVars()

			if botToken != tt.expectedBotToken {
				t.Errorf("expected botToken %q, got %q", tt.expectedBotToken, botToken)
			}
			if guildID != tt.expectedGuildID {
				t.Errorf("expected guildID %q, got %q", tt.expectedGuildID, guildID)
			}
			if (err != nil && err.Error() != tt.expectedError) || (err == nil && tt.expectedError != "") {
				t.Errorf("expected error %q, got %v", tt.expectedError, err)
			}
		})
	}
}

// contains checks if a substring is present in a string.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}
