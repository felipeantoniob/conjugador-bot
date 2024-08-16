package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type TenseMood struct {
	Mood  string `json:"mood"`
	Tense string `json:"tense"`
}

type TenseMoodChoice struct {
	Name  string    `json:"name"`
	Value TenseMood `json:"value"`
}

const (
	errTenseNameNotFound = "Tense name not found"
)

var TenseMoodChoices = []TenseMoodChoice{
	{
		Name: "Present",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Preterite",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Pretérito",
		},
	},
	{
		Name: "Imperfect",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Imperfecto",
		},
	},
	{
		Name: "Conditional",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Condicional",
		},
	},
	{
		Name: "Future",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Futuro",
		},
	},
	{
		Name: "Present perfect",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Preterite perfect (Past anterior)",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Pretérito anterior",
		},
	},
	{
		Name: "Pluperfect (Past perfect)",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Pluscuamperfecto",
		},
	},
	{
		Name: "Conditional perfect",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Condicional perfecto",
		},
	},
	{
		Name: "Future perfect",
		Value: TenseMood{
			Mood:  "Indicativo",
			Tense: "Futuro perfecto",
		},
	},
	{
		Name: "Present subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Presente",
		},
	},
	{
		Name: "Imperfect subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Imperfecto",
		},
	},
	{
		Name: "Future subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Futuro",
		},
	},
	{
		Name: "Present perfect subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Presente perfecto",
		},
	},
	{
		Name: "Pluperfect (Past perfect) subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Pluscuamperfecto",
		},
	},
	{
		Name: "Future perfect subjunctive",
		Value: TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Pretérito anterior",
		},
	},
	{
		Name: "Imperative",
		Value: TenseMood{
			Mood:  "Imperativo Afirmativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Negative Imperative",
		Value: TenseMood{
			Mood:  "Imperativo Negativo",
			Tense: "Presente",
		},
	},
}

func GetTenseMoodChoices() []*discordgo.ApplicationCommandOptionChoice {
	TenseMoodChoicesWithNameAsValue := make([]*discordgo.ApplicationCommandOptionChoice, len(TenseMoodChoices))

	for i, choice := range TenseMoodChoices {
		TenseMoodChoicesWithNameAsValue[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  choice.Name,
			Value: choice.Name,
		}
	}

	return TenseMoodChoicesWithNameAsValue
}

func GetValueByName(name string) (TenseMood, error) {
	choicesMap := make(map[string]TenseMood)
	for _, choice := range TenseMoodChoices {
		choicesMap[choice.Name] = choice.Value
	}

	value, ok := choicesMap[name]
	if !ok {
		return TenseMood{}, fmt.Errorf("%s: %s", errTenseNameNotFound, name)
	}

	return value, nil
}
