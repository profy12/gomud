package main

import (
	"errors"
	"testing"
)

func TestCreatePlayer(t *testing.T) {
	p, err := playerLoad("test_joueur")
	if !errors.Is(err, ErrUnregistered) {
		t.Fatalf("Le player de test ne devrait pas exister: %v", err)
	}
	if p.HpMax != 10 {
		t.Fatal("New player should have 10 hp")
	}
	err = p.playerSave()
	if err != nil {
		t.Fatalf("Saving test player to file: %v", err)
	}
	err = p.del()
	if err != nil {
		t.Fatalf("Unable to delete test user: %v", err)
	}
}

func TestMovePlayer(t *testing.T) {
	p, err := playerLoad("test_joueur")
	if !errors.Is(err, ErrUnregistered) {
		t.Fatalf("Le player de test ne devrait pas exister: %v", err)
	}
	RoomLoad(p.RoomId)
	if !p.isExit("u") {
		t.Fatalf("up devrait être une sortie de la piece par defaut")
	}
	if p.isExit("z") {
		t.Fatalf("z ne sera jamais une sortie")
	}
	p.del()
}
