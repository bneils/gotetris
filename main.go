package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	cells          [][]Cell
	cellsWidth     int
	cellsHeight    int
	sprite         Sprite
	frameNumber    uint64
	averageFps     RollingAverageFPS
	piece          *Piece
	gameTick       uint64
	lastMoveTick   uint64
	keyboard       Keyboard
	upcomingPieces [3]Piece
}

type Cell struct {
	rgba  color.RGBA
	solid bool
}

func MakeGame() *Game {
	const CELLS_WIDTH int = 10
	const CELLS_HEIGHT int = 20

	game := Game{
		cells:        make([][]Cell, 20),
		cellsWidth:   CELLS_WIDTH,
		cellsHeight:  CELLS_HEIGHT,
		frameNumber:  0,
		averageFps:   MakeFPSAverage(30),
		piece:        nil,
		gameTick:     0,
		lastMoveTick: 0,
		keyboard:     MakeKeyboard(),
	}

	// Load sprite
	data, err := base64.StdEncoding.DecodeString(BLOCK_IMAGE_B64)
	if err != nil {
		log.Fatal(err)
	}
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	game.sprite = MakeSprite(game.cellsWidth, game.cellsHeight, src)

	// Initialize cells
	for y := range game.cellsHeight {
		game.cells[y] = make([]Cell, game.cellsWidth)
	}

	// Initialize pieces
	game.upcomingPieces = [3]Piece{MakePiece(), MakePiece(), MakePiece()}
	piece := MakePiece()
	game.piece = &piece

	return &game
}

func (g *Game) Update() error {
	g.gameTick++

	if g.piece == nil {
		return nil
	}

	if g.keyboard.KeyPulse(ebiten.KeyUp, g.gameTick, 20) {
		// Rotate
		prev_state := g.piece.cells
		g.piece.cells = Rotate(g.piece, Clockwise)
		collision, shift := g.CheckOverlap()
		// Decide whether to rollback?
		if collision {
			if shift == 0 {
				g.piece.cells = prev_state
			} else {
				g.piece.x += shift
				collision2, _ := g.CheckOverlap()
				if collision2 {
					g.piece.x -= shift
					g.piece.cells = prev_state
				}
			}
		}
	}

	moveDown := false
	moveX := 0

	if g.keyboard.KeyPulse(ebiten.KeyDown, g.gameTick, 10) {
		// Move down
		moveDown = true
	}
	if g.keyboard.KeyPulse(ebiten.KeyLeft, g.gameTick, 10) {
		// Shift left
		moveX = -1
	}
	if g.keyboard.KeyPulse(ebiten.KeyRight, g.gameTick, 10) {
		// Shift right
		moveX = 1
	}

	if g.gameTick-g.lastMoveTick >= 60 {
		moveDown = true
	}

	if moveX != 0 {
		g.piece.x += moveX
		if collision, _ := g.CheckOverlap(); collision {
			g.piece.x -= moveX
		}
	}

	if moveDown {
		g.lastMoveTick = g.gameTick
		g.piece.y++
		if collision, _ := g.CheckOverlap(); collision {
			g.piece.y--
			// Make piece solid
			for ly, row := range g.piece.cells {
				for lx, cell := range row {
					if cell.solid {
						gy := g.piece.y + ly
						if gy < 0 {
							// Game over
							println("Game over")
							g.piece = nil
							return nil
						}
						gx := g.piece.x + lx
						g.cells[gy][gx] = cell
					}
				}
			}

			// Check if rows solid, then shift down
			shifts := [20]int{}
			for y, row := range g.cells {
				row_solid := true
				for _, cell := range row {
					if !cell.solid {
						row_solid = false
					}
				}
				if row_solid {
					for prevY := range y {
						shifts[prevY]++
					}
					for x := range row {
						row[x] = Cell{solid: false}
					}
				}
			}

			for y := len(shifts) - 2; y > 0; y-- {
				if shifts[y] != 0 {
					g.cells[y+shifts[y]] = g.cells[y]
					g.cells[y] = make([]Cell, 10)
				}
			}

			g.TakeNextPiece()
		}
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Gotetris?")

	game := MakeGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
