package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

type Player struct {
	Pseudo         string `yaml:"pseudo"`
	State          string
	DiscordPseudo  string `yaml:"discord_pseudo"`
	DiscordChannel string `yaml:"discord_channel"`
	Description    string
	ScoreId        string `yaml:"score_message_id"`
	RoomScreenId   string `yaml:"room_screen_id"`
	RoomId         string `yaml:"room_id"`
	HpMax          uint
	MpMax          uint
	HpCur          uint
	MpCur          uint
	Tired          uint
	// Session   Session
}

func (p *Player) tick() {
	if p.State != "active" {
		return
	}
	RoomLoad(defaultRoom)
	changed := false
	if p.HpCur < p.HpMax {
		p.HpCur++
		changed = true
	}
	if p.MpCur < p.MpMax {
		p.MpCur++
		changed = true
	}
	if p.Tired > 0 {
		p.Tired--
		changed = true
	}
	if changed {
		log.Printf("%v regen", p.Pseudo)
		//p.msg("Vous vous régénérez")
		go p.refreshTopic()
		p.score(false)
	}
}

func (p *Player) mv(exit string) error {
	ex := rooms[p.RoomId].Exits[exit]
	oldRoom, _ := RoomLoad(p.RoomId)
	newRoom, err := RoomLoad(ex.Target)
	if err != nil {
		return fmt.Errorf("%v moving from %v to %v: %v", p.Pseudo, rooms[p.RoomId].Name, ex.Target, err)
	}
	delete(oldRoom.Positions, p.DiscordPseudo)
	p.RoomId = ex.Target
	newRoom.Positions[p.DiscordPseudo] = &Position{ArrivedAt: time.Now()}
	p.look(false)
	p.Tired++
	p.score(false)
	p.playerSave()
	return nil
}

func (p Player) msg(msg string) {
	dg.ChannelMessageSend(p.DiscordChannel, msg)
}

func (p Player) getFileName() string {
	return fmt.Sprintf("%s/%s.yml", playerDataDir, p.DiscordPseudo)
}
func (p Player) del() error {
	err := os.Remove(p.getFileName())
	if err != nil {
		log.Printf("Deleting user: %v", err)
		return err
	}
	delete(players, p.DiscordPseudo)
	return nil
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
func (p *Player) look(create bool) {
	embed := rooms[p.RoomId].Display()
	var message *discordgo.Message
	var err error
	if p.RoomScreenId != "" && !create {
		message, err = dg.ChannelMessageEditEmbed(p.DiscordChannel, p.RoomScreenId, &embed)
	} else {
		message, err = dg.ChannelMessageSendEmbed(p.DiscordChannel, &embed)
	}
	if err != nil {
		log.Printf("Erreur pendant l'envoie de l'embed : %v", err)
	}
	p.RoomScreenId = message.ID
	p.playerSave()
}
func (p *Player) score(create bool) {
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
			{
				Name:   "Fatigue",
				Value:  fmt.Sprintf("%d", p.Tired),
				Inline: true,
			},
		},
	}
	var message *discordgo.Message
	var err error
	if p.ScoreId != "" && !create {
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

func (p Player) isExit(ex string) bool {
	for d := range rooms[p.RoomId].Exits {
		if ex == d {
			return true
		}
	}
	return false
}
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
	p, exists := players[id]
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
			RoomId:        defaultRoom,
			DiscordPseudo: id,
		}
		players[id] = pl
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
	players[id] = p
	_, err = RoomLoad(p.RoomId)
	if err != nil {
		log.Printf("command look, unable to load room: %v", err)
		p.msg("Vous n'êtes nulle part et ce n'est pas normal !")
	}
	return p, nil
}
