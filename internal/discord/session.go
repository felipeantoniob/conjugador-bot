package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	errBotInit          = "error initializing bot"
	errClosingSession   = "Error closing Discord session: %v"
	errDiscordWSOpen    = "error opening websocket connection to Discord"
	errRegisterCommands = "failed to register commands"
)

// SessionFactory is an interface for creating new Discord sessions.
type SessionFactory interface {
	New(token string) (Session, error)
}

// DefaultSessionFactory is the default implementation of SessionFactory.
type DefaultSessionFactory struct{}

// New creates a new Discord session with the provided token.
func (f *DefaultSessionFactory) New(token string) (Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &DiscordSession{session}, nil
}

// Session represents a Discord session and provides methods for interacting with it.
type Session interface {
	Open() error
	Close() error
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error)
	AddHandler(handler interface{}) func()
	GetUserID() string
	SetIntents(intents discordgo.Intent)
}

// DiscordSession wraps a discordgo.Session and provides additional functionality.
type DiscordSession struct {
	*discordgo.Session
}

// Open starts the Discord session.
func (ds *DiscordSession) Open() error {
	return ds.Session.Open()
}

// Close shuts down the Discord session.
func (ds *DiscordSession) Close() error {
	return ds.Session.Close()
}

// ApplicationCommandCreate creates a new application command for the specified application and guild.
func (ds *DiscordSession) ApplicationCommandCreate(applicationID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error) {
	return ds.Session.ApplicationCommandCreate(applicationID, guildID, cmd)
}

// AddHandler registers a handler function for events.
func (ds *DiscordSession) AddHandler(handler interface{}) func() {
	return ds.Session.AddHandler(handler)
}

// GetUserID returns the ID of the user associated with the session.
func (ds *DiscordSession) GetUserID() string {
	return ds.State.User.ID
}

func (ds *DiscordSession) SetIntents(intent discordgo.Intent) {
	ds.Identify.Intents = intent
}

// CreateSession initializes a new Discord session with the provided factory and token.
func CreateSession(factory SessionFactory, token string) (Session, error) {
	session, err := factory.New(token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBotInit, err)
	}
	session.SetIntents(discordgo.IntentsGuildMessages)
	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("%s: %w", errDiscordWSOpen, err)
	}
	return session, nil
}

// CloseSession gracefully closes the given Discord session and logs any error.
func CloseSession(s Session) {
	if err := s.Close(); err != nil {
		log.Printf(errClosingSession, err)
	}
}
