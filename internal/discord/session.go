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

// sessionInterface abstracts the methods used from discordgo.Session
type sessionInterface interface {
	Open() error
	Close() error
	SetIntents(intents discordgo.Intent)
	AddHandler(handler interface{}) func()
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error)
}

// Implement sessionInterface for discordgo.Session
type discordSession struct {
	*discordgo.Session
}

func (d *discordSession) Open() error {
	return d.Session.Open()
}

func (d *discordSession) Close() error {
	return d.Session.Close()
}

func (d *discordSession) SetIntents(intents discordgo.Intent) {
	d.Identify.Intents = intents
}

// NewSession creates a new discordSession based on the provided token
func NewSession(token string) (sessionInterface, error) {
	sess, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	return &discordSession{sess}, nil
}

// CreateSession creates a new Discord session with the provided bot token.
func CreateSession(token string, sessionFactory func(string) (sessionInterface, error)) (sessionInterface, error) {
	session, err := sessionFactory("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBotInit, err)
	}

	session.SetIntents(discordgo.IntentsGuildMessages)

	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("%s: %w", errDiscordWSOpen, err)
	}

	return session, nil
}

// CloseSession closes the Discord session and logs any error.
func CloseSession(session sessionInterface) {
	if err := session.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	}
}

// RegisterHandlersAndCommands registers event handlers and commands with the Discord session.
func RegisterHandlersAndCommands(session sessionInterface, guildID string) error {
	session.AddHandler(onReady)
	if err := registerCommands(session, guildID); err != nil {
		return fmt.Errorf("%s: %w", errRegisterCommands, err)
	}
	return nil
}
