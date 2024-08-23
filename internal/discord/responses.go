package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/internal/db"
)

const (
	errSendingResponse      = "Failed to send interaction response"
	errSendingErrorResponse = "Failed to send error interaction response"
)

func createConjugationEmbed(infinitive string, verb *db.Verb) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s - %s", infinitive, db.NullStringToString(verb.VerbEnglish)),
		Color: 16711807,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Tiempo", Value: verb.Tense},
			{Name: "Modo", Value: verb.Mood},
			{Name: "yo", Value: db.NullStringToString(verb.Form1s), Inline: true},
			{Name: "tú", Value: db.NullStringToString(verb.Form2s), Inline: true},
			{Name: "él/ella/Ud.", Value: db.NullStringToString(verb.Form3s), Inline: true},
			{Name: "nosotros", Value: db.NullStringToString(verb.Form1p), Inline: true},
			{Name: "vosotros", Value: db.NullStringToString(verb.Form2p), Inline: true},
			{Name: "ellos/ellas/Uds.", Value: db.NullStringToString(verb.Form3p), Inline: true},
		},
	}
}

// sendInteractionResponse sends a response to the interaction with the provided data
func sendInteractionResponse(session *discordgo.Session, interaction *discordgo.Interaction, responseData *discordgo.InteractionResponseData) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	}

	if err := session.InteractionRespond(interaction, response); err != nil {
		fmt.Printf("%s: %v\n", errSendingResponse, err)
	}
}

// sendConjugationResponse sends a response with the provided embed message
func sendConjugationResponse(session *discordgo.Session, interaction *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	responseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	sendInteractionResponse(session, interaction, responseData)
}

// sendErrorInteractionResponse sends an error message as a response to a Discord interaction
func sendErrorInteractionResponse(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	responseData := &discordgo.InteractionResponseData{
		Content: errorMessage,
	}
	sendInteractionResponse(session, interaction, responseData)
}
