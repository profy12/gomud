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

type Session struct {
	s      *discordgo.Session
	m      *discordgo.MessageCreate
	player Player
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatalf("Le token n'est pas défini dans le fichier .env")
	}
	fmt.Println("Démarrage du bot")
	dg, err := discordgo.New("Bot " + token)
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

	fmt.Println("message reçu")
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
	switch pl.State {
	// Player is new, we ask him his pseudo
	case "needPseudo":
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		s.ChannelMessageSend(m.ChannelID, "Bienvenue, peux tu me donner ton pseudo ?")
		pl.State = "waitPseudo"
	// Player should send us a pseudo
	case "waitPseudo":
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Est tu certain que tu souhaite %s comme pseudonyme ? (o/n)", m.Content))
		pl.State = "confirmPseudo"
		pl.Pseudo = m.Content
	case "confirmPseudo":
		if m.Content == "o" {
			log.Printf("%s state is %s", m.Author.Username, pl.State)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bienvenue %s !", pl.Pseudo))
			pl.State = "active"
		} else {
			pl.State = "needPseudo"
		}
	case "active":
		log.Printf("%s state is %s", m.Author.Username, pl.State)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hey %v je crois qu'il serait temps d'implémenter des commandes, veux tu m'aider ?", pl.Pseudo))
	default:
		log.Printf("%s state -%s- is unknown", m.Author.Username, pl.State)
	}
	//	if err != nil {
	//		log.Printf("Impossible de charger le joueur : %v", err)
	//		if errors.Is(err, ErrUnregistered) {
	//			playerCreate(s, m)
	//		} else {
	//			s.ChannelMessageSend(m.ChannelID, "Impossible de charger ton profil")
	//			log.Printf("Unable to load profil of %s", m.Author.Username)
	//		}
	//	} else {
	//
	//		log.Printf("Welcome back %v", player[m.Author.Username].Pseudo)
	//	}
}
