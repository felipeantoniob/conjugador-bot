package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/database"
	"github.com/felipeantoniob/goConjugationBot/db"
	"github.com/felipeantoniob/goConjugationBot/utils"
)

const (
	errTenseData = "Error getting tense data"
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	} else {
		log.Printf("No handler found for command: %s", i.ApplicationCommandData().Name)
	}
}

func HandleConjugate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	database, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := database.Query(
		"SELECT * FROM verbs WHERE infinitive = ? AND mood = ? AND tense = ?",
		infinitive, tenseMoodObject.Mood, tenseMoodObject.Tense,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verb db.Verb
	if rows.Next() {
		if err := rows.Scan(
			&verb.Infinitive, &verb.Mood, &verb.Tense, &verb.VerbEnglish,
			&verb.Form1s, &verb.Form2s, &verb.Form3s,
			&verb.Form1p, &verb.Form2p, &verb.Form3p,
		); err != nil {
			return nil, err
		}
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

func sendConjugationResponse(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func SendErrorInteractionResponse(s *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errorMessage,
		},
	}
	s.InteractionRespond(interaction, response)
}
