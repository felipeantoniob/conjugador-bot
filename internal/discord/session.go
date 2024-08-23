package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	errBotInit          = "error initializing bot"
	errDiscordWSOpen    = "error opening websocket connection to Discord"
	errRegisterCommands = "failed to register commands"
)

// CreateSession creates a new Discord session with the provided bot token.
func CreateSession(token string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBotInit, err)
	}
	session.Identify.Intents = discordgo.IntentsGuildMessages
	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("%s: %w", errDiscordWSOpen, err)
	}
	return session, nil
}

// CloseSession closes the Discord session and logs any error.
func CloseSession(session *discordgo.Session) {
	if err := session.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	}
}
