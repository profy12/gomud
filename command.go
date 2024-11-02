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
					p.msg("parler, p, gossip : Parler permet d'envoyer un message à tout le monde")
				case "save":
					p.msg("Permet de sauvegarder ton personnage, pas très utile car doit se faire tout seul")
				case "score":
					p.msg("score, sc : Permet d'actualiser ou de créer un nouveau message de score, score pour actualiser et score add pour créer un nouveau message de score")
				case "look":
					p.msg("look, l : Permet de regarder autour de vous")
				default:
					p.msg("Commande inconnue, ou non documentée")
				}
			} else {
				p.msg("Mauvaise syntaxe, utilise help <commande>")
			}
		} else {
			p.msg("Commandes dispos : look, score, parler, save")
		}
	case "score", "sc":
		if option {
			if options == "add" {
				p.score(true)
			} else {
				p.msg("score option unknown see help score")
			}
		}
		p.score(false)
	case "look", "l":
		r, err := RoomLoad(p.RoomId)
		if err != nil {
			log.Printf("command look, unable to load room: %v", err)
			p.msg("Vous n'êtes nulle part et ce n'est pas normal !")
		} else {
			p.msg(fmt.Sprintf("Vous êtes dans %v", r.Name))
		}
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
