package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// TenseMood represents a grammatical mood and tense.
type TenseMood struct {
	Mood  string `json:"mood"`
	Tense string `json:"tense"`
}

// TenseMoodChoice represents a choice for a tense mood.
type TenseMoodChoice struct {
	Name  string    `json:"name"`
	Value TenseMood `json:"value"`
}

const (
	errTenseNameNotFound = "Tense name not found"
)

// tenseMoodChoices holds the available tense mood choices.
var tenseMoodChoices = []TenseMoodChoice{
	{"Present", TenseMood{"Indicativo", "Presente"}},
	{"Preterite", TenseMood{"Indicativo", "Pretérito"}},
	{"Imperfect", TenseMood{"Indicativo", "Imperfecto"}},
	{"Conditional", TenseMood{"Indicativo", "Condicional"}},
	{"Future", TenseMood{"Indicativo", "Futuro"}},
	{"Present perfect", TenseMood{"Indicativo", "Presente"}},
	{"Preterite perfect (Past anterior)", TenseMood{"Indicativo", "Pretérito anterior"}},
	{"Pluperfect (Past perfect)", TenseMood{"Indicativo", "Pluscuamperfecto"}},
	{"Conditional perfect", TenseMood{"Indicativo", "Condicional perfecto"}},
	{"Future perfect", TenseMood{"Indicativo", "Futuro perfecto"}},
	{"Present subjunctive", TenseMood{"Subjuntivo", "Presente"}},
	{"Imperfect subjunctive", TenseMood{"Subjuntivo", "Imperfecto"}},
	{"Future subjunctive", TenseMood{"Subjuntivo", "Futuro"}},
	{"Present perfect subjunctive", TenseMood{"Subjuntivo", "Presente perfecto"}},
	{"Pluperfect (Past perfect) subjunctive", TenseMood{"Subjuntivo", "Pluscuamperfecto"}},
	{"Future perfect subjunctive", TenseMood{"Subjuntivo", "Pretérito anterior"}},
	{"Imperative", TenseMood{"Imperativo Afirmativo", "Presente"}},
	{"Negative Imperative", TenseMood{"Imperativo Negativo", "Presente"}},
}

// tenseMoodMap provides a quick lookup for tense moods by name.
var tenseMoodMap = createTenseMoodMap()

func createTenseMoodMap() map[string]TenseMood {
	m := make(map[string]TenseMood)
	for _, choice := range tenseMoodChoices {
		m[choice.Name] = choice.Value
	}
	return m
}

// getTenseMoodChoices returns a list of discordgo.ApplicationCommandOptionChoice based on TenseMoodChoices.
func getTenseMoodChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(tenseMoodChoices))
	for i, choice := range tenseMoodChoices {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  choice.Name,
			Value: choice.Name,
		}
	}
	return choices
}

// getValueByName returns the TenseMood for the given name, or an error if not found.
func getValueByName(name string) (TenseMood, error) {
	value, ok := tenseMoodMap[name]
	if !ok {
		return TenseMood{}, fmt.Errorf("%s: %s", errTenseNameNotFound, name)
	}
	return value, nil
}
