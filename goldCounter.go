package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

type goldCounter struct {
	gold  int
	x, y  int
	color color.RGBA
}

func (counter *goldCounter) drawCurrentGoldText(screen *ebiten.Image, textOps *text.DrawOptions, face font.Face) {
	textFace := text.NewGoXFace(face)
	textOps.GeoM.Translate(float64(counter.x), float64(counter.y))
	textOps.ColorScale.ScaleWithColor(counter.color)
	text.Draw(screen, "Gold: "+strconv.Itoa(counter.gold), textFace, textOps)
	textOps.ColorScale.Reset()
	textOps.GeoM.Reset()
}
