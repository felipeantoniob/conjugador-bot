package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/commands"
	"github.com/felipeantoniob/goConjugationBot/config"
	"github.com/felipeantoniob/goConjugationBot/database"
	"github.com/felipeantoniob/goConjugationBot/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	errEnvLoad           = "error loading env variables"
	ErrBotInit           = "Error initializing bot"
	ErrCmdCreate         = "Cannot create command"
	ErrTenseData         = "Error getting tense data"
	ErrDBQuery           = "Error querying database"
	ErrDBScan            = "Error scanning database row"
	ErrDiscordWSOpen     = "Error opening websocket connection to Discord"
	ErrTenseNameNotFound = "Tense name not found"
)

type Verb struct {
	Infinitive  string
	Mood        string
	Tense       string
	VerbEnglish string
	Form1s      string
	Form2s      string
	Form3s      string
	Form1p      string
	Form2p      string
	Form3p      string
}

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
	if err := database.InitDB(); err != nil {
		return err
	}
	defer database.CloseDB()

	// Create a new Discord session
	dgSession, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBotInit, err)
	}

	dgSession.Identify.Intents = discordgo.IntentsGuildMessages

	if err := dgSession.Open(); err != nil {
		return fmt.Errorf("%s: %w", ErrDiscordWSOpen, err)
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
