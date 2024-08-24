package main

import (
	"fmt"
	"log"
	"os"

	"github.com/felipeantoniob/goConjugationBot/internal/db"
	"github.com/felipeantoniob/goConjugationBot/internal/discord"
	"github.com/felipeantoniob/goConjugationBot/internal/env"
	u "github.com/felipeantoniob/goConjugationBot/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	errEnvLoad          = "error loading env variables"
	errBotInit          = "Error initializing bot"
	errDiscordWSOpen    = "Error opening websocket connection to Discord"
	errRegisterCommands = "failed to register commands"
	errDBInit           = "failed to initialize database"
	errDBClose          = "error closing database: %v"
	errRetrieveEnvVars  = "failed to retrieve environment variables"

	msgBotRunning = "Bot is now running. Press CTRL-C to exit."
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	if err := env.LoadEnv(&env.GodotenvLoader{}); err != nil {
		return fmt.Errorf("%s: %w", errEnvLoad, err)
	}

	botToken, guildID, err := env.GetRequiredEnvVars()
	if err != nil {
		return fmt.Errorf("%s: %w", errRetrieveEnvVars, err)
	}

	if err := db.InitDB("./internal/db/verbs.db"); err != nil {
		return fmt.Errorf("%s: %w", errDBInit, err)
	}
	defer closeDatabase()

	session, err := discord.CreateSession(&discord.DefaultSessionFactory{}, botToken)
	if err != nil {
		return fmt.Errorf("%s: %w", errBotInit, err)
	}
	defer discord.CloseSession(session)

	if err := discord.SetupCommands(session, guildID, discord.CommandRegistry); err != nil {
		return fmt.Errorf("%s: %w", errRegisterCommands, err)
	}

	fmt.Println(msgBotRunning)

	sigCh := make(chan os.Signal, 1)
	sig := u.WaitForShutdown(sigCh)
	fmt.Printf("Received shutdown signal: %v\n", sig)

	return nil
}

func closeDatabase() {
	if err := db.CloseDB(); err != nil {
		log.Printf(errDBClose, err)
	}
}
