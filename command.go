package main

import (
	"fmt"
	"log"
	"strings"
)

func (p *Player) parseCommand(msg string) error {
	//log.Printf("Commande reçu avant split : %v", msg)
	command, options, option := strings.Cut(msg, " ")
	log.Printf("Commande %s envoyée par %v", command, p.Pseudo)
	switch command {
	case "help", "h":
		if option {
			subcommand, _, opt := strings.Cut(options, " ")
			if !opt {
				p.msg(fmt.Sprintf("Aide de %v", subcommand))
				switch subcommand {
				case "parler":
					p.msg("Parler permet d'envoyer un message à tout le monde, tu peux aussi utiliser p ou gossip pour aller plus vite !")
				case "save":
					p.msg("Permet de sauvegarder ton personnage, pas très utile car doit se faire tout seul")
				default:
					p.msg("Commande inconnue, ou non documentée")
				}
			} else {
				p.msg("Mauvaise syntaxe, utilise help <commande>")
			}
		} else {
			p.msg("Commandes dispos : parler, save")
		}
	case "score", "sc":
		p.score()
	case "tick":
		p.tick()
	case "desc", "description":
		if option {
			p.Description = options
		}
	case "save":
		p.playerSave()
	case "parler", "p", "gossip", "g":
		if option {
			gossip(options, p.Pseudo)
		} else {
			p.msg("mais que veux tu dire exactement ?")
		}
	default:
		p.msg(fmt.Sprintf("%v: commande inconnue", command))
	}
	return nil
}
