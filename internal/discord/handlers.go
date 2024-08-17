package discord

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/goConjugationBot/internal/db"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"conjugate": handleConjugate,
}

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
		sendErrorInteractionResponse(s, i.Interaction, "Infinitive or tense not provided.")
		return
	}

	tenseMoodObject, err := getValueByName(tense)
	if err != nil {
		log.Println(errTenseData, err)
		sendErrorInteractionResponse(s, i.Interaction, "Error getting tense data.")
		return
	}

	verb, err := fetchVerbFromDB(infinitive, tenseMoodObject)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorInteractionResponse(s, i.Interaction, "Verb not found.")
			return
		}

		log.Println("Error fetching verb:", err)
		sendErrorInteractionResponse(s, i.Interaction, "Error querying database.")
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
