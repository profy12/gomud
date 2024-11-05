package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var (
	room = make(map[string]*Room)
)

type Room struct {
	Name  string `yaml:"name"`
	Desc  string `yaml:"desc"`
	Id    string `yaml:"id"`
	Temp  int
	Peace bool
	Dark  bool
	Exits map[string]Exit
}

type Exit struct {
	Target string
	Lock   bool
	Desc   string
	Hidden bool
}

func (r Room) Display() discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       r.Name,
		Description: r.Desc,
		Color:       0x00ff00,
	}
	exitLabel := "<"
	for ex, exit := range r.Exits {
		if exit.Hidden {
			continue
		}
		f := discordgo.MessageEmbedField{
			Name:   ex,
			Value:  exit.Desc,
			Inline: false,
		}
		exitLabel = exitLabel + " " + ex
		embed.Fields = append(embed.Fields, &f)
	}
	exitLabel = exitLabel + " >"
	//append(embed.Fields, )
	embed.Footer = &discordgo.MessageEmbedFooter{Text: exitLabel}
	return embed

}

func RoomLoad(id string) (*Room, error) {
	var r *Room
	r, exists := room[id]
	if !exists {
		filename := fmt.Sprintf("%s/%s.yml", roomDataDir, id)
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("Loading room: %v", err)
			return nil, err
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Loading room: %v", err)
		}
		yaml.Unmarshal(data, &r)
		room[id] = r
	}
	return r, nil
}
