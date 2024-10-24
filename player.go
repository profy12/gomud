package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var (
	players = []Player{}
	player  = make(map[string]*Player)
)

type Player struct {
	Pseudo string `yaml:"pseudo"`
	//State     string
	//Hitpoint  int
	//Manapoint int
	// Session   Session
}

//func (p Player, msg string) Msg() {
//	p.Session.s.ChannelMessageSend(p.Session.m.ChannelID, msg)
//}

//	func (p Player) Score() {
//		fmt.Println("Nom : ", p.Pseudo)
//		fmt.Println("Points de vie :", p.Hitpoint)
//	}
func playerCreate(s *discordgo.Session, m *discordgo.MessageCreate) (*Player, error) {
	s.ChannelMessageSendReply(m.ChannelID, "Quel est ton nom ?", m.MessageReference)
	return nil, nil
}

func playerLoad(id string) error {
	var p *Player
	p, exists := player[id]
	if exists {
		return nil
	}
	fileName := fmt.Sprintf("%s/%s.yml", playerDataDir, id)
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	//pl := Player{}
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	player[id] = p
	return nil
}
