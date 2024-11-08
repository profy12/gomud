package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	players          = make(map[string]*Player)
	topicBusy       = make(map[string]bool)
	ErrUnregistered = errors.New("unregistered: the player is unknown")
)

// send message to all players
func gossip(msg string, from string) {
	message := fmt.Sprintf("(%v): %v", from, msg)
	for _, pl := range players {
		pl.msg(message)
	}
}

func tick() {
	for {
		//log.Printf("Wait for next tick")
		time.Sleep(time.Minute)
		log.Printf("Start a new tick")
		for _, pl := range players {
			log.Printf("Tick for %v", pl.Pseudo)
			go pl.tick()
		}
	}
}
