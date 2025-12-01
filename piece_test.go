package main

import (
	"testing"
)

func TestCheckOverlapWithCells(t *testing.T) {
	game := MakeGame()
	game.cells[19][4] = Cell{solid: true}
	game.piece = &Piece{}
}
