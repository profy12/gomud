package main

import (
	"fmt"
)

var (
	commands []Command
)
type Command struct {
	Name string
	State bool
	Hitpoint   int
	Manapoint   int
	Session Session
}

func (p Player) Score() {
	fmt.Println("Nom : ", p.Name)
	fmt.Println("Points de vie : ", p.hp)
}
