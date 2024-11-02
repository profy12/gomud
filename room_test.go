package main

import "testing"

func TestRoomLoad(t *testing.T) {
	r, err := RoomLoad("taverne")
	if err != nil {
		t.Fatalf("Impossible de charger la pièce de la taverne : %v", err)
	}
	if r.Name != "La taverne" {
		t.Fatal("S'il n'y a pas de taverne alors on a soucis")
	}
}
