package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	players         []*Player
	player          = make(map[string]*Player)
	ErrUnregistered = errors.New("unregistered: the player is unknown")
)

type Player struct {
	Pseudo         string `yaml:"pseudo"`
	State          string
	DiscordPseudo  string `yaml:"discord_pseudo"`
	DiscordChannel string `yaml:"discord_channel"`
	//Hitpoint  int
	//Manapoint int
	// Session   Session
}

func (p Player) msg(msg string) {
	dg.ChannelMessageSend(p.DiscordChannel, msg)
}

//	func (p Player) Score() {
//		fmt.Println("Nom : ", p.Pseudo)
//		fmt.Println("Points de vie :", p.Hitpoint)
//	}
func (p Player) playerSave() error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s/%s.yml", playerDataDir, p.DiscordPseudo)
	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		return err
	}
	//time.Sleep(time.Minute)
	p.msg("You have been saved to file")
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
