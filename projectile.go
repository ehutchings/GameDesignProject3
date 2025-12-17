package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type effectType int

const (
	damageOnly = effectType(iota)
	slow
	stun
)

type effect struct {
	typeOfEffect effectType
	duration     int
	strength     int
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
	isHitscan              bool
}

type projectileManager struct {
	projectiles []projectile
}

func (projectile *projectile) Update(enemies []*enemy) {
	projectile.updatePosition()
	if projectile.isOnTarget() {
		if projectile.AreaOfEffectRadius != 0 {
			aoeCollider := resolv.NewCircle(float64(projectile.x-projectile.sprite.Bounds().Dx()/2),
				float64(projectile.y-projectile.sprite.Bounds().Dy()/2), projectile.AreaOfEffectRadius)
			for _, enemy := range enemies {
				if enemy != nil && aoeCollider.DistanceTo(enemy.collider) <= projectile.AreaOfEffectRadius {
					enemy.health -= projectile.inheritedDamage
				}
			}
		} else if projectile.targetEnemy != nil { //Enemy has already been defeated by another projectile
			projectile.targetEnemy.health -= projectile.inheritedDamage
		}
	}
}

func (projectile *projectile) updatePosition() {
	if projectile.x < projectile.targetX {
		projectile.x += projectile.speed
	} else if projectile.x > projectile.targetX {
		projectile.x -= projectile.speed
	}
	if projectile.y < projectile.targetY {
		projectile.y += projectile.speed
	} else if projectile.y > projectile.targetY {
		projectile.y -= projectile.speed
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
	drawOps.GeoM.Translate(float64(projectile.x), float64(projectile.y))
	screen.DrawImage(projectile.sprite, drawOps)
	drawOps.GeoM.Reset()
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
