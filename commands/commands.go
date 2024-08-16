package commands

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
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
				Choices:     GetTenseMoodChoices(),
			},
		},
	},
}
