package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	ErrCmdCreate = "Cannot create command"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"conjugate": HandleConjugate,
}

func RegisterCommands(s *discordgo.Session, guildID string) error {
	if err := RegisterBotCommands(s, guildID); err != nil {
		return fmt.Errorf("failed to register bot commands: %w", err)
	}

	s.AddHandler(HandleInteraction)
	return nil
}

// RegisterBotCommands registers the commands with Discord and logs any errors.
func RegisterBotCommands(s *discordgo.Session, guildID string) error {
	for _, cmd := range Commands {
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd); err != nil {
			return fmt.Errorf("%s '%v' command: %w", ErrCmdCreate, cmd.Name, err)
		}
	}
	return nil
}
