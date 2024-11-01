package main

import (
	"fmt"
	"log"
	"strings"
)

var (
	commands []Command
)

type Command struct {
	Name  string
	State bool
}

func (p Player) parseCommand(msg string) error {
	log.Printf("Commande reçu avant split : %v", msg)
	command, options, option := strings.Cut(msg, " ")
	log.Printf("Commande reçu : %v", command)
	switch command {
	case "help", "h":
		if option {
			p.msg(fmt.Sprintf("Tu as besoin d'aide sur %v", options))
		} else {
			p.msg("Commandes dispos : parler, save")
		}
	case "save":
		p.playerSave()
	case "parler", "p", "gossip":
		if option {
			gossip(options, p.Pseudo)
		} else {
			p.msg("mais que veux tu dire exactement ?")
		}
	default:
		p.msg(fmt.Sprintf("%v: command unknown", command))
	}
	return nil
}
