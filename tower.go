package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type towerType int

const (
	crossbow = towerType(iota)
	blackHole
	infernalEye
)

type tower struct {
	typeOfTower               towerType
	spritesheet               *ebiten.Image
	x, y                      int
	baseDamage                int
	rangeRadius               int
	baseCostMod               float64
	firing                    bool
	frameLength, currentFrame int
	rangeCollider             *resolv.Circle
	targetEnemy               *enemy
}

func (tower *tower) Draw(drawOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	drawOps.GeoM.Translate(float64(tower.x), float64(tower.y))
	frame := tower.currentFrame * TILE_WIDTH
	screen.DrawImage(tower.spritesheet.SubImage(image.Rect(frame, 0,
		frame+TILE_WIDTH, TILE_HEIGHT)).(*ebiten.Image), drawOps)
	drawOps.GeoM.Reset()
}

func (tower *tower) getTarget(enemies []*enemy) {
	tower.targetEnemy = nil
	highestDistanceTravelled := 0
	for _, currentEnemy := range enemies {
		if tower.rangeCollider.IsIntersecting(currentEnemy.collider) ||
			currentEnemy.collider.IsContainedBy(tower.rangeCollider) {
			if currentEnemy.distanceTravelled > highestDistanceTravelled {
				tower.targetEnemy = currentEnemy
			}
		}
	}
}

func newCrossbowTower(x, y int) *tower {
	sheet := LoadEmbeddedImage("Towers", "crossbowSpriteSheet.png")
	radius := 250
	return &tower{
		typeOfTower:   crossbow,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    1,
		baseCostMod:   1,
		rangeRadius:   radius,
		firing:        false,
		currentFrame:  0,
		frameLength:   3,
		rangeCollider: resolv.NewCircle(float64(x-TILE_WIDTH/2), float64(y-TILE_HEIGHT/2), float64(radius)),
	}
}

func (tower *tower) fireProjectile(projManager *projectileManager) {
	if tower.typeOfTower == crossbow {
		sprite := LoadEmbeddedImage("Projectiles", "crossbowBolt.png")
		newProjectile := projectile{
			x:               tower.x,
			y:               tower.y,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           8,
			effect:          damageOnly,
		}
		projManager.projectiles = append(projManager.projectiles, &newProjectile)
	}
}
