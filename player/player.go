package main

import (
	"fmt"
)

var (
	players = []string{}
)
type Player struct {
	Name string
	state string
	hp   int
	mp   int
}

func (p Player) Score() {
	fmt.Println("Nom : ", p.Name)
	fmt.Println("Points de vie : ", p.hp)
}
