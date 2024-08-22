package main

import (
	"fmt"
	"log"

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

	msgBotRunning       = "Bot is now running. Press CTRL-C to exit."
	msgShutdownReceived = "Shutdown signal received, exiting."
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	if err := env.LoadEnv(); err != nil {
		return u.WrapError(errEnvLoad, err)
	}

	botToken, guildID, err := env.GetRequiredEnvVars()
	if err != nil {
		return u.WrapError(errRetrieveEnvVars, err)
	}

	if err := db.InitDB("./internal/db/verbs.db"); err != nil {
		return u.WrapError(errDBInit, err)
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

	fmt.Println(msgBotRunning)
	u.WaitForShutdown()
	fmt.Println(msgShutdownReceived)

	return nil
}

func closeDatabase() {
	if err := db.CloseDB(); err != nil {
		log.Printf(errDBClose, err)
	}
}
