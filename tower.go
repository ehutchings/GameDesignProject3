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
	snowflake
)

type tower struct {
	typeOfTower               towerType
	spritesheet               *ebiten.Image
	x, y                      int
	baseDamage                int
	rangeRadius               float64
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
			tower.fireProjectile(projectileManager, tower.targetEnemy.x+TILE_WIDTH/2, tower.targetEnemy.y+TILE_WIDTH/2)
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

func newCrossbowTower(x, y, damageMod int) tower {
	sheet := LoadEmbeddedImage("Towers", "crossbowSpriteSheet.png")
	radius := 300.0
	return tower{
		typeOfTower:   crossbow,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    12 * damageMod,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   60,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   3,
		rangeCollider: resolv.NewCircle(float64(x+TILE_WIDTH/2), float64(y+TILE_HEIGHT/2), radius),
	}
}

func newVoidLauncherTower(x, y, damageMod int) tower {
	sheet := LoadEmbeddedImage("Towers", "voidLauncherSpriteSheet.png")
	radius := 250.0
	return tower{
		typeOfTower:   voidLauncher,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    8 * damageMod,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   120,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   4,
		rangeCollider: resolv.NewCircle(float64(x+TILE_WIDTH/2), float64(y+TILE_HEIGHT/2), radius),
	}
}

func newInfernalEyeTower(x, y, damageMod int) tower {
	sheet := LoadEmbeddedImage("Towers", "infernalEyeSpriteSheet.png")
	radius := 350.0
	return tower{
		typeOfTower:   infernalEye,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    2 * damageMod,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   20,
		cooldown:      0,
		currentFrame:  0,
		frameLength:   4,
		rangeCollider: resolv.NewCircle(float64(x+TILE_WIDTH/2), float64(y+TILE_HEIGHT/2), radius),
	}
}

func newSnowflakeTower(x, y, damageMod int) tower {
	sheet := LoadEmbeddedImage("Towers", "snowflakeSpriteSheet.png")
	radius := 150.0
	return tower{
		typeOfTower:   snowflake,
		spritesheet:   sheet,
		x:             x,
		y:             y,
		baseDamage:    2 * damageMod,
		rangeRadius:   radius,
		firing:        false,
		firingDelay:   30,
		cooldown:      0,
		frameLength:   6,
		currentFrame:  0,
		rangeCollider: resolv.NewCircle(float64(x+TILE_WIDTH/2), float64(y+TILE_HEIGHT/2), radius),
	}
}

func (tower *tower) fireProjectile(projManager *projectileManager, targetX, targetY int) {
	if tower.typeOfTower == crossbow {
		sprite := LoadEmbeddedImage("Projectiles", "crossbowBolt.png")
		newProjectile := projectile{
			x:               tower.x + TILE_WIDTH/2,
			y:               tower.y + TILE_HEIGHT/2,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           14,
			effect:          nil,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
	if tower.typeOfTower == voidLauncher {
		sprite := LoadEmbeddedImage("Projectiles", "voidSphere.png")
		newProjectile := projectile{
			x:               tower.x + TILE_WIDTH/2,
			y:               tower.y + TILE_HEIGHT/2,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           10,
			effect: &effect{
				typeOfEffect: stun,
				duration:     30,
				strength:     0,
			},
			AreaOfEffectRadius: 128,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
	if tower.typeOfTower == infernalEye {
		sprite := LoadEmbeddedImage("Projectiles", "infernalBeam.png")
		newProjectile := projectile{
			x:               tower.x + TILE_WIDTH/2,
			y:               tower.y + TILE_HEIGHT/2,
			targetEnemy:     tower.targetEnemy,
			targetX:         targetX,
			targetY:         targetY,
			sprite:          sprite,
			xDirection:      1,
			yDirection:      1,
			inheritedDamage: tower.baseDamage,
			speed:           15,
			effect: &effect{
				typeOfEffect: burn,
				duration:     60,
				interval:     20,
				strength:     1,
			},
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
	if tower.typeOfTower == snowflake {
		newProjectile := projectile{
			x:               tower.x + TILE_WIDTH/2,
			y:               tower.y + TILE_HEIGHT/2,
			targetEnemy:     tower.targetEnemy,
			targetX:         tower.x + TILE_WIDTH/2,
			targetY:         tower.y + TILE_HEIGHT/2,
			sprite:          nil,
			xDirection:      0,
			yDirection:      0,
			inheritedDamage: tower.baseDamage,
			speed:           0,
			effect: &effect{
				typeOfEffect:    slow,
				duration:        60,
				interval:        0,
				durationElapsed: 0,
				strength:        2,
			},
			AreaOfEffectRadius: tower.rangeRadius,
		}
		projManager.projectiles = append(projManager.projectiles, newProjectile)
	}
}
