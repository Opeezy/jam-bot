package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	guildID         = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	botToken        = flag.String("token", "", "Bot access token")
	removeCommands  = flag.Bool("rmcmd", true, "Remove all commands after shutdown or not")
	spotifyClientId = flag.String("spid", "", "The Spotify client ID")
	spofifySecret   = flag.String("spsecret", "", "The Spotify client secret")
	environment     = flag.String("env", "", "The current app environment")

	infoLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "register-spotify",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"register-spotify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
					Title:    "Register Spotify Account",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "Username",
									Label:       "Enter your username",
									Style:       discordgo.TextInputShort,
									Placeholder: "",
									Required:    true,
									MaxLength:   300,
									MinLength:   1,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "Password",
									Label:       "Enter your password",
									Style:       discordgo.TextInputShort,
									Placeholder: "",
									Required:    true,
									MaxLength:   300,
									MinLength:   1,
								},
							},
						},
					},
				},
			})

			if err != nil {
				log.Printf("error at slash command: basic-command. Reason: %s", err)
			}
		},
	}
)

func init() { flag.Parse() }

func main() {
	log.Println("creating loggers...")
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	infoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("environment=%s", *environment)

	infoLog.Println("initializing client...")
	discord, err := discordgo.New("Bot " + *botToken)
	if err != nil {
		errorLog.Fatalf("client failed to initialize: %s", err)
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		case discordgo.InteractionModalSubmit:

			// TODO: add logic to register a user's Spotify account into a database.

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Registration complete.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			if err != nil {
				errorLog.Printf("unable to respond to interaction with user: %s-%s", i.Member.User.Username, i.Member.User.ID)
				return
			}

		}
	})

	discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		infoLog.Printf("logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)
	})

	infoLog.Println("opening connection...")

	err = discord.Open()
	if err != nil {
		errorLog.Fatalf("unable to open connection: %s", err)
	}
	defer discord.Close()

	infoLog.Println("adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, *guildID, v)
		if err != nil {
			errorLog.Fatalf("cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	infoLog.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	infoLog.Println("terminating program...")

	// Removing all commands added previously if RemoveCommands==true
	if *removeCommands {
		infoLog.Println("removing commands...")
		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, *guildID, v.ID)
			if err != nil {
				errorLog.Panicf("cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}

func satisfySongRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	infoLog.Printf("incoming message from %s", m.Author.Username)

	// To prevent the bot messages triggering this function
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Recieving the incoming message from the user
}
