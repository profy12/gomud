package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	players = []Player{}
	player  = make(map[string]*Player)
)

type Player struct {
	DiscordId string
	Pseudo    string
	State     string
	Hitpoint  int
	Manapoint int
	// Session   Session
}

//func (p Player, msg string) Msg() {
//	p.Session.s.ChannelMessageSend(p.Session.m.ChannelID, msg)
//}

func (p Player) Score() {
	fmt.Println("Nom : ", p.Pseudo)
	fmt.Println("Points de vie :", p.Hitpoint)
}
func loadPlayer(id string) (*Player, error) {
	p, exists := player[id]
	if exists {
		return p, nil
	}
	*p = Player{}
	fileName := fmt.Sprintf("%s/%s.yml", playerDataDir, id)
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
