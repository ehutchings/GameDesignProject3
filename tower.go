package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type tower struct {
	spritesheet               *ebiten.Image
	x, y                      int
	baseDamage                int
	rangeRadius               int
	baseCostMod               float64
	firing                    bool
	frameLength, currentFrame int
}

func (tower *tower) Draw(drawOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	drawOps.GeoM.Translate(float64(tower.x), float64(tower.y))
	frame := tower.currentFrame * TILE_WIDTH
	screen.DrawImage(tower.spritesheet.SubImage(image.Rect(frame, 0,
		frame+TILE_WIDTH, TILE_HEIGHT)).(*ebiten.Image), drawOps)
	drawOps.GeoM.Reset()
}

func newCrossbowTower(x, y int) *tower {
	sheet := LoadEmbeddedImage("Towers", "crossbowSpriteSheet.png")
	return &tower{
		spritesheet:  sheet,
		x:            x,
		y:            y,
		baseDamage:   1,
		baseCostMod:  1,
		rangeRadius:  250,
		firing:       false,
		currentFrame: 0,
		frameLength:  3,
	}
}
