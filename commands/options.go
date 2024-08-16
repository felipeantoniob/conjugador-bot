package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/felipeantoniob/goConjugationBot/models"
)

const (
	errTenseNameNotFound = "Tense name not found"
)

var TenseMoodChoices = []models.TenseMoodChoice{
	{
		Name: "Present",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Preterite",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Pretérito",
		},
	},
	{
		Name: "Imperfect",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Imperfecto",
		},
	},
	{
		Name: "Conditional",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Condicional",
		},
	},
	{
		Name: "Future",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Futuro",
		},
	},
	{
		Name: "Present perfect",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Preterite perfect (Past anterior)",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Pretérito anterior",
		},
	},
	{
		Name: "Pluperfect (Past perfect)",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Pluscuamperfecto",
		},
	},
	{
		Name: "Conditional perfect",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Condicional perfecto",
		},
	},
	{
		Name: "Future perfect",
		Value: models.TenseMood{
			Mood:  "Indicativo",
			Tense: "Futuro perfecto",
		},
	},
	{
		Name: "Present subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Presente",
		},
	},
	{
		Name: "Imperfect subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Imperfecto",
		},
	},
	{
		Name: "Future subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Futuro",
		},
	},
	{
		Name: "Present perfect subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Presente perfecto",
		},
	},
	{
		Name: "Pluperfect (Past perfect) subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Pluscuamperfecto",
		},
	},
	{
		Name: "Future perfect subjunctive",
		Value: models.TenseMood{
			Mood:  "Subjuntivo",
			Tense: "Pretérito anterior",
		},
	},
	{
		Name: "Imperative",
		Value: models.TenseMood{
			Mood:  "Imperativo Afirmativo",
			Tense: "Presente",
		},
	},
	{
		Name: "Negative Imperative",
		Value: models.TenseMood{
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

func GetValueByName(name string) (models.TenseMood, error) {
	choicesMap := make(map[string]models.TenseMood)
	for _, choice := range TenseMoodChoices {
		choicesMap[choice.Name] = choice.Value
	}

	value, ok := choicesMap[name]
	if !ok {
		return models.TenseMood{}, fmt.Errorf("%s: %s", errTenseNameNotFound, name)
	}

	return value, nil
}
