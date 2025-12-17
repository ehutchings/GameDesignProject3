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
	targetX, targetY       int
	sprite                 *ebiten.Image
	xDirection, yDirection float64
	inheritedDamage        int
	speed                  int
	effect                 effectType
	targetEnemy            *enemy
}

type projectileManager struct {
	projectiles []*projectile
}

func (projectile *projectile) Update() {
	projectile.updatePosition()
	if projectile.x == projectile.targetX && projectile.y == projectile.targetY {
		//make enemies take damage
	}
}

func (projectile *projectile) updatePosition() {
	if projectile.x < projectile.targetX {
		projectile.x += projectile.speed
	} else if projectile.x > projectile.targetX {
		projectile.x = projectile.targetX
	}
	if projectile.y < projectile.targetY {
		projectile.y += projectile.speed
	} else if projectile.y > projectile.targetY {
		projectile.y = projectile.targetY
	}
}

func (projectile *projectile) Draw(screen *ebiten.Image) {
	//TODO
}

func (projectileManager *projectileManager) DrawProjectiles(screen *ebiten.Image) {
	//TODO
}

func (projectileManager *projectileManager) UpdateProjectiles() {
	for index, currentProjectile := range projectileManager.projectiles {
		currentProjectile.Update()
		if currentProjectile.x == currentProjectile.targetX && currentProjectile.y == currentProjectile.targetY {
			projectileManager.removeProjectileAtIndex(index)
		}
	}
}

func (projectileManager *projectileManager) removeProjectileAtIndex(index int) {
	if len(projectileManager.projectiles) >= 2 {
		projectileManager.projectiles = append(projectileManager.projectiles[:index],
			projectileManager.projectiles[index+1:]...)
	} else {
		projectileManager.projectiles = projectileManager.projectiles[:0]
	}
}
