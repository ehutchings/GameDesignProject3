package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type towerType int

const (
	crossbow = towerType(iota)
	voidLauncher
	infernalEye
)

type damageType int

const (
	firedProjectile = damageType(iota)
	AreaOfEffect
)

type tower struct {
	typeOfTower               towerType
	spritesheet               *ebiten.Image
	x, y                      int
	baseDamage                int
	rangeRadius               float64
	baseCostMod               float64
	firing                    bool
	firingDelay, cooldown     int
	frameLength, currentFrame int
	rangeCollider             *resolv.Circle
	targetEnemy               *enemy
}

func (tower *tower) Update(enemies []*enemy, projectileManager *projectileManager) {
	tower.getTarget(enemies)
	if tower.targetEnemy != nil {
		tower.firing = true
		tower.cooldown -= 1
		if tower.cooldown%(tower.firingDelay/tower.frameLength) == 0 {
			tower.currentFrame += 1
		}
		if tower.cooldown <= 0 {
			tower.fireProjectile(projectileManager, tower.targetEnemy.x, tower.targetEnemy.y)
			tower.cooldown = tower.firingDelay
			tower.currentFrame = 0
		}
	} else {
		tower.firing = false
		tower.currentFrame = 0
	}
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
		if tower.rangeCollider.DistanceTo(currentEnemy.collider) < tower.rangeRadius {
			if currentEnemy.distanceTravelled > highestDistanceTravelled {
				tower.targetEnemy = currentEnemy
			}
		}
	}
}

func newCrossbowTower(x, y int) tower {
	sheet := LoadEmbeddedImage("Towers", "crossbowSpriteSheet.png")
	radius := 300.0
	return tower{
		typeOfTower:   crossbow,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    3,
		baseCostMod:   1,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   60,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   3,
		rangeCollider: resolv.NewCircle(float64(x-TILE_WIDTH/2), float64(y-TILE_HEIGHT/2), float64(radius)),
	}
}

func newVoidLauncherTower(x, y int) tower {
	sheet := LoadEmbeddedImage("Towers", "voidLauncherSpriteSheet.png")
	radius := 200.0
	return tower{
		typeOfTower:   voidLauncher,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    2,
		baseCostMod:   2,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   120,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   4,
		rangeCollider: resolv.NewCircle(float64(x-TILE_WIDTH/2), float64(y-TILE_HEIGHT/2), float64(radius)),
	}
}

func newInfernalEyeTower(x, y int) tower {
	sheet := LoadEmbeddedImage("Towers", "infernalEyeSpriteSheet.png")
	radius := 350.0
	return tower{
		typeOfTower:   infernalEye,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    1,
		baseCostMod:   1.5,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   10,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   4,
		rangeCollider: resolv.NewCircle(float64(x-TILE_WIDTH/2), float64(y-TILE_HEIGHT/2), float64(radius)),
	}
}

func (tower *tower) fireProjectile(projManager *projectileManager, targetX, targetY int) {
	if tower.typeOfTower == crossbow {
		sprite := LoadEmbeddedImage("Projectiles", "crossbowBolt.png")
		newProjectile := projectile{
			x:               tower.x,
			y:               tower.y,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           14,
			isHitscan:       false,
			effect:          nil,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
	if tower.typeOfTower == voidLauncher {
		sprite := LoadEmbeddedImage("Projectiles", "voidSphere.png")
		newProjectile := projectile{
			x:               tower.x,
			y:               tower.y,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           8,
			isHitscan:       false,
			effect: &effect{
				typeOfEffect: stun,
				duration:     60,
				strength:     0,
			},
			AreaOfEffectRadius: 60,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
	if tower.typeOfTower == infernalEye {
		sprite := LoadEmbeddedImage("Projectiles", "infernalBeam.png")
		newProjectile := projectile{
			x:               tower.x,
			y:               tower.y,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           0,
			isHitscan:       true,
			effect:          nil,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
}
