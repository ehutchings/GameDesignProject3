package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/paths"
)

type enemySpawn struct {
	sprite *ebiten.Image
	x, y   int
}

type enemy struct {
	spritesheet            *ebiten.Image
	x, y                   int
	xDirection, yDirection int
	path                   *paths.Path
}

func newEnemySpawn(x, y int) *enemySpawn {
	image := LoadEmbeddedImage("", "enemySpawn.png")
	return &enemySpawn{
		sprite: image,
		x:      x,
		y:      y,
	}
}
