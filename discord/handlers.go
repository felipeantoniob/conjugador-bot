package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/db"
	"github.com/felipeantoniob/goConjugationBot/utils"
)

const (
	errTenseData            = "Error getting tense data"
	errSendingResponse      = "Failed to send interaction response"
	errSendingErrorResponse = "Failed to send error interaction response"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	} else {
		log.Printf("No handler found for command: %s", i.ApplicationCommandData().Name)
	}
}

func handleConjugate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := makeOptionMap(options)

	infinitive, tense, err := extractInfinitiveAndTense(optionMap)
	if err != nil {
		log.Println("Missing required options:", err)
		SendErrorInteractionResponse(s, i.Interaction, "Infinitive or tense not provided.")
		return
	}

	tenseMoodObject, err := GetValueByName(tense)
	if err != nil {
		log.Println(errTenseData, err)
		SendErrorInteractionResponse(s, i.Interaction, "Error getting tense data.")
		return
	}

	verb, err := fetchVerbFromDB(infinitive, tenseMoodObject)
	if err != nil {
		log.Println("Error fetching verb:", err)
		SendErrorInteractionResponse(s, i.Interaction, "Error querying database.")
		return
	}

	conjugationEmbed := createConjugationEmbed(infinitive, verb)
	sendConjugationResponse(s, i.Interaction, conjugationEmbed)
}

func makeOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func extractInfinitiveAndTense(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption) (infinitive string, tense string, err error) {
	if opt, exists := optionMap["infinitive"]; exists {
		infinitive = opt.StringValue()
	} else {
		err = fmt.Errorf("infinitive not found")
	}

	if opt, exists := optionMap["tense"]; exists {
		tense = opt.StringValue()
	} else {
		err = fmt.Errorf("tense not found")
	}

	return infinitive, tense, err
}

func fetchVerbFromDB(infinitive string, tenseMoodObject TenseMood) (*db.Verb, error) {
	ctx := context.Background()
	sqlDB, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	queries := db.New(sqlDB)

	verb, err := queries.GetVerbByInfinitiveMoodTense(ctx, db.GetVerbByInfinitiveMoodTenseParams{Infinitive: infinitive, Mood: tenseMoodObject.Mood, Tense: tenseMoodObject.Tense})
	if err != nil {
		return nil, err
	}

	return &verb, nil
}

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
