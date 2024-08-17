package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/internal/utils"
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
		return nil, utils.WrapError(errBotInit, err)
	}
	session.Identify.Intents = discordgo.IntentsGuildMessages
	if err := session.Open(); err != nil {
		return nil, utils.WrapError(errDiscordWSOpen, err)
	}
	return session, nil
}

// CloseSession closes the Discord session and logs any error.
func CloseSession(session *discordgo.Session) {
	if err := session.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	}
}

// RegisterHandlersAndCommands registers event handlers and commands with the Discord session.
func RegisterHandlersAndCommands(session *discordgo.Session, guildID string) error {
	session.AddHandler(onReady)
	if err := registerCommands(session, guildID); err != nil {
		return utils.WrapError(errRegisterCommands, err)
	}
	return nil
}
