package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ErrEnvLoad           = "Error loading env variables"
	ErrBotInit           = "Error initializing bot"
	ErrDBOpen            = "Error opening database connection"
	ErrCmdCreate         = "Cannot create command"
	ErrTenseData         = "Error getting tense data"
	ErrDBQuery           = "Error querying database"
	ErrDBScan            = "Error scanning database row"
	ErrMissingVars       = "BOT_TOKEN and GUILD_ID environment variables are required"
	ErrDiscordWSOpen     = "Error opening websocket connection to Discord"
	ErrTenseNameNotFound = "Tense name not found"
)

var (
	db *sql.DB

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "conjugate",
			Description: "Provides conjugation details for a given Spanish verb.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "infinitive",
					Description: "Verb to look up.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tense",
					Description: "Tense and mood of the chosen verb.",
					Required:    true,
					Choices:     getChoices(),
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"conjugate": handleConjugate,
	}
)

type Verb struct {
	Infinitive  string
	Mood        string
	Tense       string
	VerbEnglish string
	Form1s      string
	Form2s      string
	Form3s      string
	Form1p      string
	Form2p      string
	Form3p      string
}

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal(ErrEnvLoad)
	}

	botToken := os.Getenv("BOT_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	if botToken == "" || guildID == "" {
		log.Fatal(ErrMissingVars)
	}

	db, err = sql.Open("sqlite3", "verbs.db")
	if err != nil {
		log.Fatalf("%s: %v", ErrDBOpen, err)
	}
	defer db.Close()

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("%s: %v", ErrBotInit, err)
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		fmt.Printf("%s: %v\n", ErrDiscordWSOpen, err)
		return
	}
	defer s.Close()

	s.AddHandler(onReady)

	registerCommands(s, guildID)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func registerCommands(s *discordgo.Session, guildID string) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("%s '%v' command: %v", ErrCmdCreate, v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
}

func handleConjugate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var infinitive, tense string
	if option, ok := optionMap["infinitive"]; ok {
		infinitive = option.StringValue()
	}
	if option, ok := optionMap["tense"]; ok {
		tense = option.StringValue()
	}

	tenseMoodObject, err := getValueByName(tense)
	if err != nil {
		log.Println(ErrTenseData, err)
		sendErrorInteractionResponse(s, i.Interaction, "Error getting tense data.")
		return
	}

	rows, err := db.Query("SELECT * FROM verbs WHERE infinitive = ? AND MOOD = ? AND tense = ?", infinitive, tenseMoodObject.Mood, tenseMoodObject.Tense)
	if err != nil {
		log.Println(ErrDBQuery, err)
		sendErrorInteractionResponse(s, i.Interaction, "Error querying database.")
		return
	}
	defer rows.Close()

	var verb Verb
	for rows.Next() {
		if err := rows.Scan(&verb.Infinitive, &verb.Mood, &verb.Tense, &verb.VerbEnglish, &verb.Form1s, &verb.Form2s, &verb.Form3s, &verb.Form1p, &verb.Form2p, &verb.Form3p); err != nil {
			log.Println(ErrDBScan, err)
			sendErrorInteractionResponse(s, i.Interaction, "Error scanning database row.")
			return
		}
	}

	conjugationEmbed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s - %s", infinitive, verb.VerbEnglish),
		Color: 16711807,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Tiempo",
				Value: verb.Tense,
			},
			{
				Name:  "Modo",
				Value: verb.Mood,
			},
			{
				Name:   "yo",
				Value:  verb.Form1s,
				Inline: true,
			},
			{
				Name:   "tú",
				Value:  verb.Form2s,
				Inline: true,
			},
			{
				Name:   "él/ella/Ud.",
				Value:  verb.Form3s,
				Inline: true,
			},
			{
				Name:   "nosotros",
				Value:  verb.Form1p,
				Inline: true,
			},
			{
				Name:   "vosotros",
				Value:  verb.Form2p,
				Inline: true,
			},
			{
				Name:   "ellos/ellas/Uds.",
				Value:  verb.Form3p,
				Inline: true,
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{conjugationEmbed},
		},
	})
}

func sendErrorInteractionResponse(s *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errorMessage,
		},
	}
	s.InteractionRespond(interaction, response)
}

type TenseMood struct {
	Mood  string `json:"mood"`
	Tense string `json:"tense"`
}

type TenseMoodChoice struct {
	Name  string    `json:"name"`
	Value TenseMood `json:"value"`
}

func getChoices() []*discordgo.ApplicationCommandOptionChoice {
	TenseMoodChoicesWithNameAsValue := make([]*discordgo.ApplicationCommandOptionChoice, len(TenseMoodChoices))

	for i, choice := range TenseMoodChoices {
		TenseMoodChoicesWithNameAsValue[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  choice.Name,
			Value: choice.Name,
		}
	}

	return TenseMoodChoicesWithNameAsValue
}

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

func getValueByName(name string) (TenseMood, error) {
	choicesMap := make(map[string]TenseMood)
	for _, choice := range TenseMoodChoices {
		choicesMap[choice.Name] = choice.Value
	}

	value, ok := choicesMap[name]
	if !ok {
		return TenseMood{}, fmt.Errorf("%s: %s", ErrTenseNameNotFound, name)
	}

	return value, nil
}
