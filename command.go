package main

import (
	"fmt"
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
	params := strings.SplitN(msg, " ", 1)
	command := params[0]
	var options string
	if len(params) > 1 {
		options = params[1]
	}
	switch command {
	case "help":
		if options != "" {
			p.msg(fmt.Sprintf("Tu as besoin d'aide sur %v", options))
		} else {
			p.msg("Commandes dispos : help, save")
		}
	case "save":
		p.playerSave()
	default:
		p.msg("Command unknown")
	}
	return nil
}
