package main

import (
	"image/color"
	"math/rand/v2"
)

type Piece struct {
	cells [][]Cell
	width int
	x, y  int
}

type PieceType int

const (
	IPiece PieceType = iota
	SqPiece
	LPiece
	RLPiece
	TPiece
	ZPiece
	RZPiece
	NumberPieces
)

type ClockDirection int

const (
	Clockwise ClockDirection = iota
	CounterClockwise
)

func Rotate(p *Piece, direction ClockDirection) [][]Cell {
	cellsCopy := make([][]Cell, len(p.cells))
	for i, row := range p.cells {
		cellsCopy[i] = make([]Cell, len(row))
	}
	for i, row := range p.cells {
		l := len(row)
		for j, cell := range row {
			var x, y int
			if direction == Clockwise {
				//l-1-y, x
				x = l - 1 - i
				y = j
			} else {
				//y, l-1-x
				x = i
				y = l - 1 - j
			}
			cellsCopy[y][x] = cell
		}
	}
	return cellsCopy
}

func MakePiece() Piece {
	ptype := PieceType(rand.IntN(int(NumberPieces)))

	colors := []color.RGBA{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
	}

	color := colors[rand.IntN(len(colors))]

	var width int
	var piece Piece
	switch ptype {
	case LPiece, RLPiece, TPiece, ZPiece, RZPiece:
		width = 3
	case IPiece:
		width = 4
	case SqPiece:
		width = 2
	}

	piece.cells = make([][]Cell, width)
	for i := range piece.cells {
		piece.cells[i] = make([]Cell, width)
	}

	cell := Cell{color, true}

	switch ptype {
	case LPiece:
		piece.cells[0][1] = cell
		piece.cells[1][1] = cell
		piece.cells[2][1] = cell
		piece.cells[2][2] = cell
	case RLPiece:
		piece.cells[0][1] = cell
		piece.cells[1][1] = cell
		piece.cells[2][1] = cell
		piece.cells[2][0] = cell
	case TPiece:
		piece.cells[0][1] = cell
		piece.cells[1][1] = cell
		piece.cells[2][1] = cell
		piece.cells[1][2] = cell
	case ZPiece:
		piece.cells[0][2] = cell
		piece.cells[1][2] = cell
		piece.cells[1][1] = cell
		piece.cells[2][1] = cell
	case RZPiece:
		piece.cells[0][0] = cell
		piece.cells[1][0] = cell
		piece.cells[1][1] = cell
		piece.cells[2][1] = cell
	case SqPiece:
		piece.cells[0][0] = cell
		piece.cells[0][1] = cell
		piece.cells[1][0] = cell
		piece.cells[1][1] = cell
	case IPiece:
		piece.cells[0][0] = cell
		piece.cells[0][1] = cell
		piece.cells[0][2] = cell
		piece.cells[0][3] = cell
	}

	mid := (10 - width) / 2
	piece.x = mid
	piece.y = -width
	piece.width = width

	return piece
}

func IAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) CheckOverlap() (bool, int) {
	collision := false
	shift := 0
	for ly, row := range g.piece.cells {
		row_collides := false
		// range of game cells that obstruct
		minBadX := g.piece.width
		maxBadX := 0
		// range of game cells that are empty
		minGoodX := g.piece.width
		maxGoodX := 0
		for lx, cell := range row {
			gx := g.piece.x + lx
			gy := g.piece.y + ly
			if gy < 0 {
				continue
			}

			pieceLayer := cell.solid
			// gameLayer is opaque if out-of-bounds or occupied.
			var gameLayer bool
			if gx < 0 || gx >= g.cellsWidth ||
				gy >= g.cellsHeight {
				gameLayer = true
			} else {
				gameLayer = g.cells[gy][gx].solid
			}

			// Intersection adjusts bad bounds
			if pieceLayer && gameLayer {
				row_collides = true
				minBadX = min(minBadX, lx)
				maxBadX = max(maxBadX, lx)
			}

			if !gameLayer {
				minGoodX = min(minGoodX, lx)
				maxGoodX = max(maxGoodX, lx)
			}
		}
		if row_collides {
			row_shift := 0
			collision = true

			// No overlap means easy shift
			if maxBadX < minGoodX {
				row_shift = minGoodX - minBadX
			} else if maxGoodX < minBadX {
				row_shift = maxGoodX - maxBadX
			} else {
				// Overlap, not an easy resolution
				return true, 0
			}

			// If shift contradicts previous shifts, then avoid
			if row_shift < 0 && 0 < shift || shift < 0 && 0 < row_shift {
				return true, 0
			}
			// Suggest largest shift
			if IAbs(row_shift) > IAbs(shift) {
				shift = row_shift
			}
		}
	}
	return collision, shift
}

func (g *Game) CheckOutOfBounds() {

}

func (g *Game) TakeNextPiece() {
	piece := g.upcomingPieces[0]
	g.piece = &piece

	for i := range len(g.upcomingPieces) - 1 {
		g.upcomingPieces[i] = g.upcomingPieces[i+1]
	}
	g.upcomingPieces[len(g.upcomingPieces)-1] = MakePiece()
}
