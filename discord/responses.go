package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/db"
	"github.com/felipeantoniob/goConjugationBot/utils"
)

const (
	errTenseData            = "Error getting tense data"
	errSendingResponse      = "Failed to send interaction response"
	errSendingErrorResponse = "Failed to send error interaction response"
)

func createConjugationEmbed(infinitive string, verb *db.Verb) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s - %s", infinitive, utils.NullStringToString(verb.VerbEnglish)),
		Color: 16711807,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Tiempo", Value: verb.Tense},
			{Name: "Modo", Value: verb.Mood},
			{Name: "yo", Value: utils.NullStringToString(verb.Form1s), Inline: true},
			{Name: "tú", Value: utils.NullStringToString(verb.Form2s), Inline: true},
			{Name: "él/ella/Ud.", Value: utils.NullStringToString(verb.Form3s), Inline: true},
			{Name: "nosotros", Value: utils.NullStringToString(verb.Form1p), Inline: true},
			{Name: "vosotros", Value: utils.NullStringToString(verb.Form2p), Inline: true},
			{Name: "ellos/ellas/Uds.", Value: utils.NullStringToString(verb.Form3p), Inline: true},
		},
	}
}

// sendConjugationResponse sends a response with the provided embed message
func sendConjugationResponse(session *discordgo.Session, interaction *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	responseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	}

	if err := session.InteractionRespond(interaction, response); err != nil {
		fmt.Printf("%s: %v\n", errSendingResponse, err)
	}
}

// SendErrorInteractionResponse sends an error message as a response to a Discord interaction
func SendErrorInteractionResponse(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	responseData := &discordgo.InteractionResponseData{
		Content: errorMessage,
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	}

	if err := session.InteractionRespond(interaction, response); err != nil {
		fmt.Printf("%s: %v\n", errSendingErrorResponse, err)
	}
}
