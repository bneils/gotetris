package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/draw"
)

type Sprite struct {
	sourceImage image.Image
	variants    map[color.Color]image.Image
	variantRect image.Rectangle
}

func GetScaledSpriteDim(cellsWidth, cellsHeight int) (scaledW, scaledH int) {
	_, h := ebiten.WindowSize()
	scaledH = int(float32(h) / float32(cellsHeight))
	scaledW = scaledH //int(float32(h) / float32(cellsWidth))
	return
}

func MakeSprite(cellsWidth, cellsHeight int, source image.Image) Sprite {
	// Scale image
	w, h := GetScaledSpriteDim(cellsWidth, cellsHeight)
	scaledSource := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(scaledSource, scaledSource.Bounds(), source, source.Bounds(), draw.Src, &draw.Options{})

	spriteColors := []color.RGBA{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
	}

	// Tint images
	variants := make(map[color.Color]image.Image, len(spriteColors))
	for _, color := range spriteColors {
		coloredImage := TintImage(scaledSource, color)
		variants[color] = coloredImage
	}

	variantRect := variants[spriteColors[0]].Bounds()

	return Sprite{source, variants, variantRect}
}
