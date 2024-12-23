package discord

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/felipeantoniob/conjugador-bot/internal/db"
)

const (
	errInfinitiveNotFound = "infinitive not found"
	errTenseNotFound      = "tense not found"
	errInfinitiveOrTense  = "Infinitive or tense not provided."
	errTenseData          = "Error getting tense data."
	errVerbNotFound       = "Verb not found."
	errQueryingDatabase   = "Error querying database."
)

func handleConjugate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := makeOptionMap(options)

	infinitive, tense, err := extractInfinitiveAndTense(optionMap)
	if err != nil {
		log.Println("Missing required options:", err)
		sendErrorInteractionResponse(&DiscordSession{s}, i.Interaction, errInfinitiveOrTense)
		return
	}

	tenseMoodObject, err := getValueByName(tense)
	if err != nil {
		sendErrorInteractionResponse(&DiscordSession{s}, i.Interaction, errTenseData)
		return
	}

	verb, err := fetchVerbFromDB(infinitive, tenseMoodObject)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorInteractionResponse(&DiscordSession{s}, i.Interaction, errVerbNotFound)
			return
		}

		log.Println("Error fetching verb:", err)
		sendErrorInteractionResponse(&DiscordSession{s}, i.Interaction, errQueryingDatabase)
		return
	}

	conjugationEmbed := createConjugationEmbed(infinitive, verb)
	sendConjugationResponse(&DiscordSession{s}, i.Interaction, conjugationEmbed)
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
		err = fmt.Errorf(errInfinitiveNotFound)
	}

	if opt, exists := optionMap["tense"]; exists {
		tense = opt.StringValue()
	} else {
		err = fmt.Errorf(errTenseNotFound)
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
