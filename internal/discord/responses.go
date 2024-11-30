package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/internal/db"
)

// createConjugationEmbed generates a Discord embed message for a verb's conjugation
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

// InteractionResponder defines an interface for sending interaction responses
type InteractionResponder interface {
	InteractionRespond(interaction *discordgo.Interaction, response *discordgo.InteractionResponse) error
}

// sendInteractionResponse sends a response to the interaction with the provided data
func sendInteractionResponse(responder InteractionResponder, interaction *discordgo.Interaction, responseData *discordgo.InteractionResponseData) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	}

	if err := responder.InteractionRespond(interaction, response); err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}

// sendConjugationResponse sends a response with the provided embed message
func sendConjugationResponse(responder InteractionResponder, interaction *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	responseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	sendInteractionResponse(responder, interaction, responseData)
}

// sendErrorInteractionResponse sends an error message as a response to a Discord interaction
func sendErrorInteractionResponse(responder InteractionResponder, interaction *discordgo.Interaction, errorMessage string) {
	responseData := &discordgo.InteractionResponseData{
		Content: errorMessage,
	}
	sendInteractionResponse(responder, interaction, responseData)
}
