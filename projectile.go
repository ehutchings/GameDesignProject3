package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type effectType int

const (
	slow = effectType(iota)
	stun
	burn
)

type effect struct {
	typeOfEffect    effectType
	duration        int
	interval        int
	durationElapsed int
	strength        int
}

type projectile struct {
	x, y                   int
	targetX, targetY       int
	sprite                 *ebiten.Image
	xDirection, yDirection float64
	inheritedDamage        int
	speed                  int
	effect                 *effect
	targetEnemy            *enemy
	AreaOfEffectRadius     float64
	rotDirection           float64
}

type projectileManager struct {
	projectiles []projectile
}

func (projectile *projectile) Update(enemies []*enemy) {
	projectile.updatePosition()
	if projectile.isOnTarget() {
		if projectile.AreaOfEffectRadius != 0 {
			var aoeCollider resolv.Circle
			if projectile.sprite != nil {
				aoeCollider = *resolv.NewCircle(float64(projectile.x+projectile.sprite.Bounds().Dx()/2),
					float64(projectile.y+projectile.sprite.Bounds().Dy()/2), projectile.AreaOfEffectRadius)
			} else {
				aoeCollider = *resolv.NewCircle(float64(projectile.x), float64(projectile.y),
					projectile.AreaOfEffectRadius)
			}
			for _, enemy := range enemies {
				if enemy != nil && aoeCollider.DistanceTo(enemy.collider) <= projectile.AreaOfEffectRadius {
					enemy.health -= projectile.inheritedDamage
					if projectile.effect != nil {
						enemy.activeEffects = append(enemy.activeEffects, projectile.effect)
					}
				}
			}
		} else if projectile.targetEnemy != nil { //Enemy has already been defeated by another projectile
			projectile.targetEnemy.health -= projectile.inheritedDamage
			if projectile.effect != nil {
				projectile.targetEnemy.activeEffects = append(projectile.targetEnemy.activeEffects, projectile.effect)
			}
		}
	}
}

func (projectile *projectile) updatePosition() {
	projectile.xDirection = float64(projectile.targetX - projectile.x)
	projectile.yDirection = float64(projectile.targetY - projectile.y)
	magnitude := math.Sqrt(math.Pow(projectile.xDirection, 2) + math.Pow(projectile.yDirection, 2))
	projectile.xDirection = projectile.xDirection / magnitude
	projectile.yDirection = projectile.yDirection / magnitude

	if projectile.x < projectile.targetX {
		projectile.x += int(float64(projectile.speed) * projectile.xDirection)
	} else if projectile.x > projectile.targetX {
		projectile.x -= int(float64(projectile.speed) * -projectile.xDirection)
	}
	if projectile.y < projectile.targetY {
		projectile.y += int(float64(projectile.speed) * projectile.yDirection)
	} else if projectile.y > projectile.targetY {
		projectile.y -= int(float64(projectile.speed) * -projectile.yDirection)
	}
}

func (projectile *projectile) isOnTarget() bool {
	if projectile.x <= projectile.targetX+projectile.speed &&
		projectile.x >= projectile.targetX-projectile.speed {
		if projectile.y <= projectile.targetY+projectile.speed &&
			projectile.y >= projectile.targetY-projectile.speed {
			return true
		}
	}
	return false
}

func (projectile *projectile) Draw(screen *ebiten.Image, drawOps *ebiten.DrawImageOptions) {
	if projectile.sprite != nil {
		drawOps.GeoM.Translate(float64(-projectile.sprite.Bounds().Dx()/2), float64(-projectile.sprite.Bounds().Dy()/2))
		drawOps.GeoM.Rotate(projectile.rotDirection)
		drawOps.GeoM.Translate(float64(projectile.x), float64(projectile.y))
		screen.DrawImage(projectile.sprite, drawOps)
		drawOps.GeoM.Reset()
	}
}

func (projectileManager *projectileManager) DrawProjectiles(screen *ebiten.Image, drawOps *ebiten.DrawImageOptions) {
	for _, projectile := range projectileManager.projectiles {
		projectile.Draw(screen, drawOps)
	}
}

func (projectileManager *projectileManager) UpdateProjectiles(enemies []*enemy) {
	for index := len(projectileManager.projectiles) - 1; index >= 0; index-- {
		projectileManager.projectiles[index].Update(enemies)
		if projectileManager.projectiles[index].isOnTarget() {
			projectileManager.removeProjectileAtIndex(index)
		}
	}
}

func (projectileManager *projectileManager) removeProjectileAtIndex(index int) {
	if len(projectileManager.projectiles) >= 2 {
		projectileManager.projectiles = append(projectileManager.projectiles[:index], projectileManager.projectiles[index+1:]...)
	} else {
		projectileManager.projectiles = projectileManager.projectiles[:0]
	}
}
