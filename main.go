package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "conjugate",
			Description: "Command for demonstrating options",
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
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "qwerty",
							Value: "qwerty",
						},
						{
							Name:  "asdf",
							Value: "asdf",
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"conjugate": handleConjugate}
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading env")
	}

	botToken := os.Getenv("BOT_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(onReady)

	s.Identify.Intents = discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer s.Close()

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
}

func handleConjugate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["infinitive"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> infinitive: %s\n"
	}

	if option, ok := optionMap["tense"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> tense: %s\n"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}
