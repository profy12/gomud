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

var (
	rooms = make(map[string]*Room)
	//people = make(map[string] map[string]*Player)
)

type Room struct {
	Name      string `yaml:"name"`
	ShortDesc string `yaml:"short_desc"`
	Desc      string `yaml:"desc"`
	Id        string `yaml:"id"`
	Temp      int
	Peace     bool
	Dark      bool
	Positions map[string]*Position
	Exits     map[string]Exit
}

// allow to know : how long a player is in room
// and list players in rooms

type Position struct {
	ArrivedAt time.Time
}

// Allow to navigate between rooms
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
	r, exists := rooms[id]
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
		log.Printf("Loading room: %v", r)
		if r.Positions == nil {
			r.Positions = make(map[string]*Position)
		}
		rooms[id] = r
	}
	return r, nil
}
