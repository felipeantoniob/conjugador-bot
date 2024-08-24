package discord

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/internal/db"
)

// Mock function for NullStringToString
func mockNullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func TestCreateConjugationEmbed(t *testing.T) {
	// Mock verb data
	verb := &db.Verb{
		VerbEnglish: sql.NullString{String: "to be", Valid: true},
		Tense:       "Present",
		Mood:        "Indicative",
		Form1s:      sql.NullString{String: "soy", Valid: true},
		Form2s:      sql.NullString{String: "eres", Valid: true},
		Form3s:      sql.NullString{String: "es", Valid: true},
		Form1p:      sql.NullString{String: "somos", Valid: true},
		Form2p:      sql.NullString{String: "sois", Valid: true},
		Form3p:      sql.NullString{String: "son", Valid: true},
	}

	// Expected embed
	expectedEmbed := &discordgo.MessageEmbed{
		Title: "be - to be",
		Color: 16711807,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Tiempo", Value: "Present"},
			{Name: "Modo", Value: "Indicative"},
			{Name: "yo", Value: "soy", Inline: true},
			{Name: "tú", Value: "eres", Inline: true},
			{Name: "él/ella/Ud.", Value: "es", Inline: true},
			{Name: "nosotros", Value: "somos", Inline: true},
			{Name: "vosotros", Value: "sois", Inline: true},
			{Name: "ellos/ellas/Uds.", Value: "son", Inline: true},
		},
	}

	// Call function
	embed := createConjugationEmbed("be", verb)

	// Compare
	if embed.Title != expectedEmbed.Title {
		t.Errorf("Expected title %s but got %s", expectedEmbed.Title, embed.Title)
	}
	if embed.Color != expectedEmbed.Color {
		t.Errorf("Expected color %d but got %d", expectedEmbed.Color, embed.Color)
	}
	for i, field := range expectedEmbed.Fields {
		if embed.Fields[i].Name != field.Name || embed.Fields[i].Value != field.Value {
			t.Errorf("Expected field %v but got %v", field, embed.Fields[i])
		}
	}
}

// Mock responder
type mockResponder struct {
	shouldFail bool
}

func (mr *mockResponder) InteractionRespond(interaction *discordgo.Interaction, response *discordgo.InteractionResponse) error {
	if mr.shouldFail {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestSendInteractionResponse(t *testing.T) {
	responder := &mockResponder{}
	interaction := &discordgo.Interaction{}
	responseData := &discordgo.InteractionResponseData{}

	// Test successful response
	sendInteractionResponse(responder, interaction, responseData)

	// Test error case
	responder.shouldFail = true
	sendInteractionResponse(responder, interaction, responseData)
}

func TestSendConjugationResponse(t *testing.T) {
	responder := &mockResponder{}
	interaction := &discordgo.Interaction{}
	embed := &discordgo.MessageEmbed{
		Title: "Test Embed",
	}

	sendConjugationResponse(responder, interaction, embed)

	// Check that the response data contains the embed
	// Note: This is a simplified test; more detailed checks can be added based on how sendInteractionResponse processes the response
}

func TestSendErrorInteractionResponse(t *testing.T) {
	responder := &mockResponder{}
	interaction := &discordgo.Interaction{}
	errorMessage := "An error occurred"

	sendErrorInteractionResponse(responder, interaction, errorMessage)

	// Check that the response data contains the error message
	// Note: This is a simplified test; more detailed checks can be added based on how sendInteractionResponse processes the response
}
