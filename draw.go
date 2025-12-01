package main

import (
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (sprite *Sprite) DrawSprite(cell Cell, px, py int, screen *ebiten.Image) {

	if !cell.solid {
		return
	}

	variant := sprite.variants[cell.rgba]

	// Skip if off-screen
	drawnRect := image.Rect(px, py, px+variant.Bounds().Dx(), py+variant.Bounds().Dy())

	if drawnRect.Intersect(screen.Bounds()).Empty() {
		return
	}

	// Draw at position
	position := ebiten.GeoM{}
	position.Translate(float64(px), float64(py))
	options := ebiten.DrawImageOptions{GeoM: position}

	ebImage := ebiten.NewImageFromImage(variant)
	screen.DrawImage(ebImage, &options)
}

func (sp *Sprite) WorldToScreenCoord(x, y int) (px, py int) {
	px = x * sp.variantRect.Dx()
	py = y * sp.variantRect.Dy()
	return
}

func ScaleRect(rect image.Rectangle, scale float32) image.Rectangle {
	ScaleInt := func(x int, s float32) int {
		return int(float32(x) * s)
	}

	return image.Rect(
		ScaleInt(rect.Min.X, scale), ScaleInt(rect.Min.Y, scale),
		ScaleInt(rect.Max.X, scale), ScaleInt(rect.Max.Y, scale),
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// FPS stopwatch
	start := time.Now()

	// Draw vertical line
	lineStartX, _ := g.sprite.WorldToScreenCoord(g.cellsWidth, 0)
	line := ebiten.NewImage(5, screen.Bounds().Dy())
	line.Fill(color.White)
	position := ebiten.GeoM{}
	position.Translate(float64(lineStartX), 0)
	options := ebiten.DrawImageOptions{GeoM: position}
	screen.DrawImage(line, &options)

	// Draw blocks
	for y, row := range g.cells {
		for x, cell := range row {
			px, py := g.sprite.WorldToScreenCoord(x, y)
			g.sprite.DrawSprite(cell, px, py, screen)
		}
	}

	// Draw falling piece
	if g.piece != nil {
		for y, row := range g.piece.cells {
			for x, cell := range row {
				px, py := g.sprite.WorldToScreenCoord(g.piece.x+x, g.piece.y+y)
				g.sprite.DrawSprite(cell, px, py, screen)
			}
		}
	}

	// Draw upcoming pieces
	topY := g.sprite.variantRect.Dy()
	midX := (screen.Bounds().Dx() + lineStartX) / 2
	for _, piece := range g.upcomingPieces {
		for y, row := range piece.cells {
			for x, cell := range row {
				w := g.sprite.variantRect.Dx()
				h := g.sprite.variantRect.Dy()
				boxW := piece.width * w
				cell_px := (midX - boxW/2) + x*w
				cell_py := topY + y*h
				g.sprite.DrawSprite(cell, cell_px, cell_py, screen)
			}
		}

		topY += (piece.width + 1) * g.sprite.variantRect.Dy()
	}

	fps := Fps(start, time.Now())
	g.averageFps.InsertFPS(fps)
	averageFps := g.averageFps.AverageFPS()
	ebitenutil.DebugPrint(screen, strconv.Itoa(int(averageFps)))

	g.frameNumber++
}
