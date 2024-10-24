package main

var (
	commands []Command
)

type Command struct {
	Name  string
	State bool
}

