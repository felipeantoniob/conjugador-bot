package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	errCmdCreate = "cannot create command '%s': %w"
)

// CommandMapping combines a Discord command with its handler function
type CommandMapping struct {
	Command *discordgo.ApplicationCommand
	Handler interface{}
}

// List of CommandMappings to be registered
var CommandRegistry = []CommandMapping{
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

// SetupCommands registers commands and their handlers with the Discord session.
func SetupCommands(s Session, guildID string, commandMappings []CommandMapping) error {
	for _, m := range commandMappings {
		if _, err := s.ApplicationCommandCreate(s.GetUserID(), guildID, m.Command); err != nil {
			return fmt.Errorf(errCmdCreate, m.Command.Name, err)
		}
		s.AddHandler(m.Handler)
	}

	return nil
}
