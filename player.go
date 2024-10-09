package main

import (
	"fmt"
)

var (
	players []Player
)
type Player struct {
	Name string
	State string
	Hitpoint   int
	Manapoint   int
	Session Session
}

func (p Player) Score() {
	fmt.Println("Nom : ", p.Name)
	fmt.Println("Points de vie : ", p.hp)
}
