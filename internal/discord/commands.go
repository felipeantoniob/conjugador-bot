package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	errCmdCreate      = "Cannot create command"
	errNoHandlerFound = "No handler found for command: %s"
)

// CommandHandlerPair combines a Discord command with its handler function
type CommandHandlerPair struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// CommandRegistry holds the application commands and their handlers.
type CommandRegistry struct {
	commands []CommandHandlerPair
}

// NewCommandRegistry initializes a new CommandRegistry.
func NewCommandRegistry(commands []CommandHandlerPair) *CommandRegistry {
	return &CommandRegistry{commands: commands}
}

// GetCommands returns a slice of application commands.
func (r *CommandRegistry) GetCommands() []*discordgo.ApplicationCommand {
	commands := make([]*discordgo.ApplicationCommand, len(r.commands))
	for i, cmd := range r.commands {
		commands[i] = cmd.Command
	}
	return commands
}

// List of ApplicationCommands and their handlers
var ApplicationCommands = []CommandHandlerPair{
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "conjugate",
			Description: "Provides conjugation details for a given Spanish verb.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "infinitive",
					Description: "Verb to look up.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tense",
					Description: "Tense and mood of the chosen verb.",
					Required:    true,
					Choices:     getTenseMoodChoices(),
				},
			},
		},
		Handler: handleConjugate,
	},
	// Add more commands and handlers here as needed
}

// commandHandlers maps command names to their handlers.
var commandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

func init() {
	for _, cmd := range ApplicationCommands {
		commandHandlers[cmd.Command.Name] = cmd.Handler
	}
}

// EventHandler is a type for functions that handle Discord events.
type EventHandler interface {
	InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type DefaultEventHandler struct {
	CommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (d *DefaultEventHandler) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdName := i.ApplicationCommandData().Name
	if handler, exists := d.CommandHandlers[cmdName]; exists {
		handler(s, i)
	} else {
		log.Printf(errNoHandlerFound, cmdName)
	}
}

// RegisterCommandsAndHandlers registers event handlers and commands with the Discord session.
func RegisterCommandsAndHandlers(s Session, guildID string, applicationCommands []CommandHandlerPair, handler EventHandler) error {
	registry := NewCommandRegistry(applicationCommands)
	cmds := registry.GetCommands()

	for _, cmd := range cmds {
		if _, err := s.ApplicationCommandCreate(s.GetUserID(), guildID, cmd); err != nil {
			return fmt.Errorf("%s '%s' command: %w", errCmdCreate, cmd.Name, err)
		}
	}

	s.AddHandler(handler.InteractionHandler)

	return nil
}
