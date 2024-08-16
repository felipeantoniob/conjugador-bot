package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/commands"
	"github.com/felipeantoniob/goConjugationBot/config"
	"github.com/felipeantoniob/goConjugationBot/db"
	"github.com/felipeantoniob/goConjugationBot/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	errEnvLoad       = "error loading env variables"
	errBotInit       = "Error initializing bot"
	errDiscordWSOpen = "Error opening websocket connection to Discord"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error running the bot: %v", err)
	}
}

func run() error {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		return fmt.Errorf("%s: %w", errEnvLoad, err)
	}

	// Retrieve required environment variables
	botToken, guildID, err1 := config.GetRequiredEnvVars()
	if err1 != nil {
		return err1
	}

	// Initialize the database
	if err := db.InitDB(); err != nil {
		return err
	}
	defer db.CloseDB()

	// Create a new Discord session
	dgSession, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("%s: %w", errBotInit, err)
	}

	dgSession.Identify.Intents = discordgo.IntentsGuildMessages

	if err := dgSession.Open(); err != nil {
		return fmt.Errorf("%s: %w", errDiscordWSOpen, err)
	}
	defer dgSession.Close()

	// Register event handlers and commands
	dgSession.AddHandler(commands.OnReady)
	commands.RegisterCommands(dgSession, guildID)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	utils.WaitForShutdown()
	fmt.Println("Shutdown signal received, exiting.")

	return nil
}
