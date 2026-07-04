package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
)

type direction int

const (
	up = direction(iota)
	down
	left
	right
)

type enemyClass int

const (
	regular = enemyClass(iota)
	fast
)

type enemySpawn struct {
	sprite        *ebiten.Image
	x, y          int
	activeEnemies []*enemy
	currentWave   *wave
	intervalCount int
	initialDelay  int
}

type enemy struct {
	spritesheet                   *ebiten.Image
	enemyClass                    enemyClass
	frameLength, currentFrame     int
	frameDelay, elapsedFrameDelay int
	x, y                          int
	baseSpeed                     int
	speed                         int
	xDirection, yDirection        int
	facingDirection               direction
	path                          *paths.Path
	collider                      *resolv.ConvexPolygon
	distanceTravelled             int
	health                        int
	activeEffects                 []*effect
	goldDropped                   int
}

func newEnemySpawn(currentWave *wave, x, y int) *enemySpawn {
	sprite := LoadEmbeddedImage("", "enemySpawn.png")
	return &enemySpawn{
		sprite:        sprite,
		x:             x,
		y:             y,
		activeEnemies: make([]*enemy, 0),
		currentWave:   currentWave,
		initialDelay:  600,
	}
}

func newRegularEnemy(x, y int) *enemy {
	spriteSheet := LoadEmbeddedImage("Enemies", "enemySpriteSheet.png")
	return &enemy{
		spritesheet:   spriteSheet,
		enemyClass:    regular,
		frameLength:   2,
		frameDelay:    16,
		x:             x,
		y:             y,
		baseSpeed:     2,
		speed:         2,
		health:        50,
		goldDropped:   4,
		xDirection:    0,
		yDirection:    0,
		collider:      resolv.NewRectangle(float64(x+TILE_WIDTH/2), float64(y+TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
		activeEffects: []*effect{},
	}
}

func newStrongRegularEnemy(x, y int) *enemy {
	spriteSheet := LoadEmbeddedImage("Enemies", "strongEnemySpriteSheet.png")
	return &enemy{
		spritesheet:   spriteSheet,
		enemyClass:    regular,
		frameLength:   2,
		frameDelay:    16,
		x:             x,
		y:             y,
		baseSpeed:     2,
		speed:         2,
		health:        200,
		goldDropped:   25,
		xDirection:    0,
		yDirection:    0,
		collider:      resolv.NewRectangle(float64(x+TILE_WIDTH/2), float64(y+TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
		activeEffects: []*effect{},
	}
}

func newStrongFastEnemy(x, y int) *enemy {
	spriteSheet := LoadEmbeddedImage("Enemies", "strongFastGoblinSpriteSheet.png")
	return &enemy{
		spritesheet:   spriteSheet,
		enemyClass:    fast,
		frameLength:   2,
		frameDelay:    12,
		x:             x,
		y:             y,
		baseSpeed:     4,
		speed:         4,
		health:        100,
		goldDropped:   30,
		xDirection:    0,
		yDirection:    0,
		collider:      resolv.NewRectangle(float64(x+TILE_WIDTH/2), float64(y+TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
		activeEffects: []*effect{},
	}
}

func newFastEnemy(x, y int) *enemy {
	spriteSheet := LoadEmbeddedImage("Enemies", "fastGoblinSpriteSheet.png")
	return &enemy{
		spritesheet:   spriteSheet,
		enemyClass:    fast,
		frameLength:   2,
		frameDelay:    12,
		x:             x,
		y:             y,
		baseSpeed:     4,
		speed:         4,
		health:        25,
		goldDropped:   10,
		xDirection:    0,
		yDirection:    0,
		collider:      resolv.NewRectangle(float64(x+TILE_WIDTH/2), float64(y+TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
		activeEffects: []*effect{},
	}
}

func (enemySpawner *enemySpawn) nextEnemyInWave(pathMap *paths.Grid, base *playerBase) {
	if enemySpawner.initialDelay > 0 {
		enemySpawner.initialDelay--
	} else {
		enemySpawner.intervalCount += 1
		if enemySpawner.intervalCount == enemySpawner.currentWave.spawnInterval {
			newEnemy := enemySpawner.currentWave.removeEnemyInFront()
			newEnemy.x = enemySpawner.x
			newEnemy.y = enemySpawner.y
			newEnemy.newPath(pathMap, enemySpawner, base)
			enemySpawner.activeEnemies = append(enemySpawner.activeEnemies, &newEnemy)
			enemySpawner.intervalCount = 0
		}
	}
}

func (enemy *enemy) newPath(pathMap *paths.Grid, enemySpawner *enemySpawn, base *playerBase) {
	startingCell := pathMap.Get(enemySpawner.x/TILE_WIDTH, enemySpawner.y/TILE_HEIGHT)
	endingCell := pathMap.Get(base.x/TILE_WIDTH, base.y/TILE_HEIGHT)
	enemy.path = pathMap.GetPathFromCells(startingCell, endingCell, false, false)
}

func canEnemyPath(game *mainGame) bool {
	startingCell := game.pathMap.Get(game.enemySpawner.x/TILE_WIDTH, game.enemySpawner.y/TILE_HEIGHT)
	endingCell := game.pathMap.Get(game.base.x/TILE_WIDTH, game.base.y/TILE_HEIGHT)
	path := game.pathMap.GetPathFromCells(startingCell, endingCell, false, false)
	if path == nil || path.Length() == 0 {
		return false
	}
	return true
}

func (enemy *enemy) Update() {
	enemy.processEffects()
	if enemy.path != nil && enemy.path.Length() > 0 {
		currentCell := enemy.path.Current()
		if math.Abs(float64(currentCell.X*TILE_WIDTH)-float64(enemy.x)) <= 2 &&
			math.Abs(float64(currentCell.Y*TILE_HEIGHT)-float64(enemy.y)) <= 2 {
			enemy.path.Advance()
		}
		enemy.xDirection = 0
		if currentCell.X*TILE_WIDTH > enemy.x {
			enemy.xDirection = 1
		} else if currentCell.X*TILE_WIDTH < enemy.x {
			enemy.xDirection = -1
		}
		enemy.yDirection = 0
		if currentCell.Y*TILE_HEIGHT > enemy.y {
			enemy.yDirection = 1
		} else if currentCell.Y*TILE_HEIGHT < enemy.y {
			enemy.yDirection = -1
		}
		enemy.x += enemy.xDirection * enemy.speed
		enemy.distanceTravelled += int(math.Abs(float64(enemy.xDirection * enemy.speed)))
		enemy.y += enemy.yDirection * enemy.speed
		enemy.distanceTravelled += int(math.Abs(float64(enemy.yDirection * enemy.speed)))
		enemy.collider.SetX(float64(enemy.x + TILE_WIDTH/2))
		enemy.collider.SetY(float64(enemy.y + TILE_HEIGHT/2))
		enemy.setDirection()
		if enemy.speed != 0 {
			enemy.elapsedFrameDelay += 1
			if enemy.elapsedFrameDelay == enemy.frameDelay {
				enemy.currentFrame += 1
				enemy.elapsedFrameDelay = 0
			}
			if enemy.currentFrame >= enemy.frameLength {
				enemy.currentFrame = 0
			}
		}
	}
}

func (enemy *enemy) setDirection() {
	switch enemy.xDirection {
	case -1:
		enemy.facingDirection = left
	case 1:
		enemy.facingDirection = right
	}
	switch enemy.yDirection {
	case -1:
		enemy.facingDirection = up
	case 1:
		enemy.facingDirection = down
	}
}

func (enemy *enemy) hasReachedBase(baseX, baseY int) bool {
	if enemy.x >= baseX && enemy.x <= baseX+TILE_WIDTH {
		if enemy.y >= baseY && enemy.y <= baseY+TILE_HEIGHT {
			return true
		}
	}
	return false
}

func (enemy *enemy) processEffects() {
	if len(enemy.activeEffects) > 0 {
		for index := len(enemy.activeEffects) - 1; index >= 0; index-- {
			currentEffect := enemy.activeEffects[index]
			if currentEffect.typeOfEffect == stun {
				enemy.speed = 0
				currentEffect.durationElapsed += 1
				if currentEffect.durationElapsed >= currentEffect.duration {
					enemy.speed = enemy.baseSpeed
					enemy.removeEffectAtIndex(index)
				}
			}
			if currentEffect.typeOfEffect == slow {
				if enemy.speed != 0 {
					enemy.speed = enemy.baseSpeed / currentEffect.strength
				}
				currentEffect.durationElapsed += 1
				if currentEffect.durationElapsed >= currentEffect.duration {
					enemy.speed = enemy.baseSpeed
					enemy.removeEffectAtIndex(index)
				}
			}
			if currentEffect.typeOfEffect == burn {
				currentEffect.durationElapsed += 1
				if currentEffect.durationElapsed%currentEffect.interval == 0 {
					enemy.health -= int(currentEffect.strength)
				}
				if currentEffect.durationElapsed >= currentEffect.duration {
					enemy.speed = enemy.baseSpeed
					enemy.removeEffectAtIndex(index)
				}
			}
		}
	}
}

func (enemy *enemy) removeEffectAtIndex(index int) {
	if len(enemy.activeEffects) >= 2 {
		enemy.activeEffects = append(enemy.activeEffects[:index], enemy.activeEffects[index+1:]...)
	} else {
		enemy.activeEffects = enemy.activeEffects[:0]
	}
}

func (enemySpawner *enemySpawn) updateEnemies(stageWaves *stageWaves, bank *goldCounter,
	pathMap *paths.Grid, base *playerBase, audioManager *audioManager) {
	enemySpawner.spawnEnemies(stageWaves, pathMap, base)
	if len(enemySpawner.activeEnemies) != 0 {
		for index := len(enemySpawner.activeEnemies) - 1; index >= 0; index-- {
			enemySpawner.activeEnemies[index].Update()
			if enemySpawner.activeEnemies[index].health <= 0 {
				bank.gold += enemySpawner.activeEnemies[index].goldDropped
				if enemySpawner.activeEnemies[index].enemyClass == regular {
					audioManager.playRegularEnemyDeathSound()
				} else if enemySpawner.activeEnemies[index].enemyClass == fast {
					audioManager.playFastEnemyDeathSound()
				}
				enemySpawner.removeEnemyAtIndex(index)
				continue
			}
			if enemySpawner.activeEnemies[index].hasReachedBase(base.x, base.y) {
				base.health -= enemySpawner.activeEnemies[index].health
				if enemySpawner.activeEnemies[index].enemyClass == regular {
					audioManager.playRegularEnemyDeathSound()
				} else if enemySpawner.activeEnemies[index].enemyClass == fast {
					audioManager.playFastEnemyDeathSound()
				}
				enemySpawner.removeEnemyAtIndex(index)
			}
		}
	}
}

func (enemySpawner *enemySpawn) spawnEnemies(stageWaves *stageWaves, pathMap *paths.Grid, base *playerBase) {
	if len(enemySpawner.currentWave.enemies) > 0 {
		enemySpawner.nextEnemyInWave(pathMap, base)
	} else if len(enemySpawner.activeEnemies) == 0 {
		if len(stageWaves.waves) != 0 {
			enemySpawner.currentWave = stageWaves.getNextWave()
		}
	}
}

func (enemy *enemy) Draw(screen *ebiten.Image, drawOps *ebiten.DrawImageOptions) {
	drawOps.GeoM.Translate(float64(enemy.x), float64(enemy.y))
	frameX := enemy.currentFrame * TILE_WIDTH
	frameY := int(enemy.facingDirection) * TILE_WIDTH
	screen.DrawImage(enemy.spritesheet.SubImage(image.Rect(frameX, frameY, frameX+TILE_WIDTH,
		frameY+TILE_HEIGHT)).(*ebiten.Image), drawOps)
	drawOps.GeoM.Reset()
}

func (enemySpawner *enemySpawn) drawEnemies(screen *ebiten.Image, drawOps *ebiten.DrawImageOptions) {
	for _, enemy := range enemySpawner.activeEnemies {
		enemy.Draw(screen, drawOps)
	}
}

func (enemySpawner *enemySpawn) removeEnemyAtIndex(index int) {
	if len(enemySpawner.activeEnemies) >= 2 {
		enemySpawner.activeEnemies = append(enemySpawner.activeEnemies[:index], enemySpawner.activeEnemies[index+1:]...)
	} else {
		enemySpawner.activeEnemies = enemySpawner.activeEnemies[:0]
	}
}

func redrawEnemyPaths(game *mainGame, enemies []*enemy) {
	for _, currentEnemy := range enemies {
		startingCell := game.pathMap.Get(currentEnemy.x/TILE_WIDTH, currentEnemy.y/TILE_HEIGHT)
		endingCell := game.pathMap.Get(game.base.x/TILE_WIDTH, game.base.y/TILE_HEIGHT)
		currentEnemy.path = game.pathMap.GetPathFromCells(startingCell, endingCell,
			false, false)
	}
}
