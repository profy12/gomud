package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	room = make(map[string]*Room)
)

type Room struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
	Id   string `yaml:"id"`
}

func (r Room) Display(){
	
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
