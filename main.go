package main

import (
	"fmt"
	"log"

	"github.com/felipeantoniob/goConjugationBot/config"
	"github.com/felipeantoniob/goConjugationBot/db"
	"github.com/felipeantoniob/goConjugationBot/discord"
	"github.com/felipeantoniob/goConjugationBot/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	errEnvLoad          = "error loading env variables"
	errBotInit          = "Error initializing bot"
	errDiscordWSOpen    = "Error opening websocket connection to Discord"
	errRegisterCommands = "failed to register commands"
	errDBInit           = "failed to initialize database"
	errRetrieveEnvVars  = "failed to retrieve environment variables"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		return utils.WrapError(errEnvLoad, err)
	}

	// Retrieve required environment variables
	botToken, guildID, err := config.GetRequiredEnvVars()
	if err != nil {
		return utils.WrapError(errRetrieveEnvVars, err)
	}

	// Initialize the database
	if err := db.InitDB("./db/verbs.db"); err != nil {
		return utils.WrapError(errDBInit, err)
	}
	defer closeDatabase()

	session, err := discord.CreateSession(botToken)
	if err != nil {
		return err
	}
	defer discord.CloseSession(session)

	if err := discord.RegisterHandlersAndCommands(session, guildID); err != nil {
		return err
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	utils.WaitForShutdown()
	fmt.Println("Shutdown signal received, exiting.")

	return nil
}

// closeDatabase closes the database connection and logs any error.
func closeDatabase() {
	if err := db.CloseDB(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
