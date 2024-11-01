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

const (
	playerDataDir = "data/players"
)

var (
	dg      *discordgo.Session
	guildId string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatalf("Le token n'est pas défini dans le fichier .env")
	}
	guildId = os.Getenv("GUILD_ID")
	if guildId == "" {
		log.Fatalln("GUILD_ID doit être défini dans le fichier .env")
	}
	fmt.Println("Démarrage du bot")
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Unable to connect on Discord: %v", err)

	}
	dg.AddHandler(messageCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	// dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	//fmt.Println("message reçu")
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
	pl, err := playerLoad(m.Author.Username)
	if err != nil {
		log.Printf("Erreur lors du chargement de l'utilisateur : %s", err)
	}
setState:
	switch pl.State {
	// Player is new, we ask him his pseudo
	case "needPseudo":
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		s.ChannelMessageSend(m.ChannelID, "Bienvenue, peux tu me donner ton pseudo ?")
		pl.State = "waitPseudo"
		pl.playerSave()
	// Player should send us a pseudo
	case "waitPseudo":
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Est tu certain que tu souhaite %s comme pseudonyme ? (o/n)", m.Content))
		pl.State = "confirmPseudo"
		pl.Pseudo = m.Content
		pl.playerSave()
	case "confirmPseudo":
		switch m.Content {
		case "o", "O", "y", "Y":
			log.Printf("%s state is %s", m.Author.Username, pl.State)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bienvenue %s !", pl.Pseudo))
			pl.State = "active"
		default:
			pl.State = "waitPseudo"
			s.ChannelMessageSend(m.ChannelID, "Et donc que veux tu comme pseudo ?")
			break setState
		}
		st, err := s.GuildChannelCreate(guildId, pl.DiscordPseudo, discordgo.ChannelTypeGuildText)
		if err != nil {
			log.Printf("Unable to create player channel: %v", err)
		}
		pl.DiscordChannel = st.ID
		pl.msg("Maintenant c'est ici que ça se passe")
		pl.playerSave()
	case "active":
		err := s.ChannelMessageDelete(pl.DiscordChannel, m.ID)
		if err != nil {
			log.Printf("Unable to delete message: %v", err)
		}
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		pl.parseCommand(m.Content)
	default:
		log.Printf("%s state -%s- is unknown", m.Author.Username, pl.State)
	}
	//pl.DiscordChannel = m.ChannelID
}
