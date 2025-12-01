package main

import (
	"image"
	"image/color"
	"log"

	"github.com/crazy3lf/colorconv"
)

func TintImage(src image.Image, rgba color.RGBA) image.Image {
	srcB := src.Bounds()
	dst := image.NewRGBA(srcB)
	h, s, scaleL := colorconv.ColorToHSL(rgba)
	scaleL *= 1.5

	// Apply tint based on brightness
	for y := srcB.Min.Y; y < srcB.Max.Y; y++ {
		for x := srcB.Min.X; x < srcB.Max.X; x++ {
			_, _, l := colorconv.ColorToHSL(src.At(x, y))
			l *= scaleL

			r8, g8, b8, err := colorconv.HSLToRGB(h, s, l)
			if err != nil {
				log.Fatal(err)
			}

			dst.Set(x, y, color.RGBA{r8, g8, b8, 255})
		}
	}

	return dst
}
