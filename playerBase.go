package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type playerBase struct {
	name     string
	health   int
	x, y     int
	sprite   *ebiten.Image
	collider *resolv.ConvexPolygon
}

func newPlayerBase(x, y int) *playerBase {
	image := LoadEmbeddedImage("", "playerBase.png")
	return &playerBase{
		name:     "",
		health:   BASE_HEALTH,
		x:        x,
		y:        y,
		sprite:   image,
		collider: resolv.NewRectangle(float64(x), float64(y), TILE_WIDTH, TILE_HEIGHT),
	}
}
