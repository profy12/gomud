package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var (
	player          = make(map[string]*Player)
	topicBusy       = make(map[string]bool)
	ErrUnregistered = errors.New("unregistered: the player is unknown")
)

// send message to all players
func gossip(msg string, from string) {
	message := fmt.Sprintf("(%v): %v", from, msg)
	for _, pl := range player {
		pl.msg(message)
	}
}

func tick() {
	for {
		log.Printf("Wait for next tick")
		time.Sleep(time.Minute)
		log.Printf("Start a new tick")
		for _, pl := range player {
			log.Printf("Tick for %v", pl.Pseudo)
			go pl.tick()
		}
	}
}

type Player struct {
	Pseudo         string `yaml:"pseudo"`
	State          string
	DiscordPseudo  string `yaml:"discord_pseudo"`
	DiscordChannel string `yaml:"discord_channel"`
	Description    string
	ScoreId        string `yaml:"score_message_id"`
	HpMax          int
	MpMax          int
	HpCur          int
	MpCur          int
	// Session   Session
}

func (p *Player) tick() {
	changed := false
	if p.HpCur < p.HpMax {
		p.HpCur++
		changed = true
	}
	if p.MpCur < p.MpMax {
		p.MpCur++
		changed = true
	}
	if changed {
		log.Printf("%v regen", p.Pseudo)
		//p.msg("Vous vous régénérez")
		go p.refreshTopic()
		p.score()
	}
}

func (p Player) msg(msg string) {
	dg.ChannelMessageSend(p.DiscordChannel, msg)
}

func (p Player) refreshTopic() {
	// si ce topic est déjà en train d'être mis à jour on annule
	if topicBusy[p.DiscordChannel] {
		return
	}
	topicBusy[p.DiscordChannel] = true
	log.Printf("topic de %v va être mis à jour", p.Pseudo)
	chEdit := discordgo.ChannelEdit{
		Topic: fmt.Sprintf("%s (%s)> HP(%d/%d) MP(%d/%d)", p.Pseudo, p.State, p.HpCur, p.HpMax, p.MpCur, p.MpMax),
	}
	_, err := dg.ChannelEdit(p.DiscordChannel, &chEdit)
	if err != nil {
		log.Printf("Problème lors de l'actualisation du topic : %v", err)
	}
	log.Printf("topic de %v mis à jour", p.Pseudo)
	topicBusy[p.DiscordChannel] = false
}
func (p *Player) score() {
	embed := discordgo.MessageEmbed{
		Title:       p.Pseudo,
		Description: p.Description,
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Etat",
				Value:  p.State,
				Inline: true,
			},
			{
				Name:   "Discord",
				Value:  p.DiscordPseudo,
				Inline: true,
			},
			{
				Name:   "Health",
				Value:  fmt.Sprintf("%d/%d", p.HpCur, p.HpMax),
				Inline: false,
			},
			{
				Name:   "Mana",
				Value:  fmt.Sprintf("%d/%d", p.MpCur, p.MpMax),
				Inline: true,
			},
		},
	}
	var message *discordgo.Message
	var err error
	if p.ScoreId != "" {
		message, err = dg.ChannelMessageEditEmbed(p.DiscordChannel, p.ScoreId, &embed)
	} else {
		message, err = dg.ChannelMessageSendEmbed(p.DiscordChannel, &embed)
	}
	if err != nil {
		log.Printf("Erreur pendant l'envoie de l'embed : %v", err)
	}
	p.ScoreId = message.ID
	p.playerSave()

}
func (p Player) msgExt(msg string) {
	//dg.ChannelMessageSend(p.DiscordChannel, msg)
	embed := discordgo.MessageEmbed{
		Title:       "Test embed",
		Description: "Pas certain que ce soit suffisant",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Champ 1",
				Value:  "Valeur du champ 1",
				Inline: true,
			},
			{
				Name:   "Champ 2",
				Value:  "Valeur du champ 2",
				Inline: true,
			},
		},
	}
	_, err := dg.ChannelMessageSendEmbed(p.DiscordChannel, &embed)
	if err != nil {
		log.Printf("Erreur pendant l'envoie de l'embed : %v", err)
	}
}

//	func (p Player) Score() {
//		fmt.Println("Nom : ", p.Pseudo)
//		fmt.Println("Points de vie :", p.Hitpoint)
//	}
func (p Player) playerSave() error {
	log.Printf("Starting saving %v to file", p.Pseudo)
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s/%s.yml", playerDataDir, p.DiscordPseudo)
	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		log.Printf("Problème lors de l'écriture du fichier : %v", err)
		return err
	}
	//time.Sleep(time.Minute)
	log.Printf("%v saved to file", p.Pseudo)
	return nil
}

// load player from file or init it
// id is the discord pseudo
func playerLoad(id string) (*Player, error) {
	var p *Player
	p, exists := player[id]
	if exists {
		return p, nil
	}
	fileName := fmt.Sprintf("%s/%s.yml", playerDataDir, id)
	f, err := os.Open(fileName)
	// if file not exist we init a new one
	if err != nil {
		pl := &Player{
			State:         "needPseudo",
			Description:   "Un joueur pas ouf",
			HpMax:         10,
			MpMax:         10,
			HpCur:         1,
			MpCur:         1,
			DiscordPseudo: id,
		}
		player[id] = pl
		return pl, ErrUnregistered
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	//pl := Player{}
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	player[id] = p
	return p, nil
}
