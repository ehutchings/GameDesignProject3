package main

import "github.com/hajimehoshi/ebiten/v2"

type effectType int

const (
	damageOnly = effectType(iota)
	slow
	pull
)

type projectile struct {
	x, y                   int
	sprite                 *ebiten.Image
	xDirection, yDirection float64
	inheritedDamage        int
	speed                  float64
	effect                 effectType
}

type projectileManager struct {
	projectiles []*projectile
}
