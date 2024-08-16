package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	errCmdCreate = "Cannot create command"
)

func registerCommands(s *discordgo.Session, guildID string) error {
	if err := registerBotCommands(s, guildID); err != nil {
		return fmt.Errorf("failed to register bot commands: %w", err)
	}

	s.AddHandler(handleInteraction)
	return nil
}

// registerBotCommands registers the commands with Discord and logs any errors.
func registerBotCommands(s *discordgo.Session, guildID string) error {
	for _, cmd := range getCommands() {
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd); err != nil {
			return fmt.Errorf("%s '%v' command: %w", errCmdCreate, cmd.Name, err)
		}
	}
	return nil
}
