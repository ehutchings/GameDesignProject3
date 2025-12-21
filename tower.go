package main

import (
	"image"
	"math"

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
	direction                 float64
}

func (tower *tower) Update(enemies []*enemy, projectileManager *projectileManager, audioManager *audioManager) {
	tower.getTarget(enemies)
	if tower.targetEnemy != nil {
		x1 := float64(tower.x + TILE_WIDTH/2)
		y1 := float64(tower.y + TILE_HEIGHT/2)
		x2 := float64(tower.targetEnemy.x + TILE_WIDTH/2)
		y2 := float64(tower.targetEnemy.y + TILE_HEIGHT/2)
		//dot := x1*x2 + y1*y2
		//det := x1*y2 - y1*x2
		tower.direction = math.Atan2(y1-y2, x1-x2) + math.Pi/2
		tower.firing = true
		tower.cooldown -= 1
		if tower.cooldown%(tower.firingDelay/tower.frameLength) == 0 {
			tower.currentFrame += 1
		}
		if tower.cooldown <= 0 {
			tower.fireProjectile(projectileManager, tower.targetEnemy.x+TILE_WIDTH/2, tower.targetEnemy.y+TILE_WIDTH/2)
			if tower.typeOfTower == voidLauncher {
				audioManager.playVoidLauncherSound()
			} else if tower.typeOfTower == snowflake {
				audioManager.playSnowflakeSound()
			} else if tower.typeOfTower == infernalEye {
				audioManager.playInfernalEyeSound()
			} else if tower.typeOfTower == crossbow {
				audioManager.playCrossbowSound()
			}
			tower.cooldown = tower.firingDelay
			tower.currentFrame = 0
		}
	} else {
		tower.firing = false
		tower.currentFrame = 0
	}
}

func (tower *tower) Draw(drawOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	if tower.typeOfTower != snowflake {
		drawOps.GeoM.Translate(-TILE_WIDTH/2, -TILE_HEIGHT/2)
		drawOps.GeoM.Rotate(tower.direction)
		drawOps.GeoM.Translate(float64(tower.x+TILE_WIDTH/2), float64(tower.y+TILE_WIDTH/2))
	} else {
		drawOps.GeoM.Translate(float64(tower.x), float64(tower.y))
	}
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
				highestDistanceTravelled = currentEnemy.distanceTravelled
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
		baseDamage:    15 * damageMod,
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
		baseDamage:    35 * damageMod,
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
		baseDamage:    3 * damageMod,
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
		baseDamage:    10 * damageMod,
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
			rotDirection:    tower.direction,
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
			rotDirection:    tower.direction,
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
				duration:     45,
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
			rotDirection:    tower.direction,
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
				duration:     90,
				interval:     15,
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
