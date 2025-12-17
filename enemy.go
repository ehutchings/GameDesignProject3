package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
)

type enemySpawn struct {
	sprite        *ebiten.Image
	x, y          int
	activeEnemies []*enemy
	currentWave   *wave
	intervalCount int
}

type enemy struct {
	spritesheet            *ebiten.Image
	x, y                   int
	speed                  int
	xDirection, yDirection int
	path                   *paths.Path
	collider               *resolv.ConvexPolygon
	distanceTravelled      int
	health                 int
	goldDropped            int
}

func newEnemySpawn(currentWave *wave, x, y int) *enemySpawn {
	image := LoadEmbeddedImage("", "enemySpawn.png")
	return &enemySpawn{
		sprite:        image,
		x:             x,
		y:             y,
		activeEnemies: make([]*enemy, 0),
		currentWave:   currentWave,
	}
}

func newEnemy(x, y, speed int) *enemy {
	image := LoadEmbeddedImage("", "enemy.png")
	return &enemy{
		spritesheet: image,
		x:           x,
		y:           y,
		speed:       speed,
		health:      10,
		goldDropped: 1,
		xDirection:  0,
		yDirection:  0,
		collider:    resolv.NewRectangle(float64(x-TILE_WIDTH/2), float64(y-TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
	}
}

func (enemySpawner *enemySpawn) nextEnemyInWave(pathMap *paths.Grid, base *playerBase) {
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
		enemy.collider.SetX(float64(enemy.x - TILE_WIDTH/2))
		enemy.collider.SetY(float64(enemy.y - TILE_HEIGHT/2))
	}
}

func (enemySpawner *enemySpawn) updateEnemies(stageWaves *stageWaves, bank *goldCounter,
	pathMap *paths.Grid, base *playerBase) {
	enemySpawner.spawnEnemies(stageWaves, pathMap, base)
	if len(enemySpawner.activeEnemies) != 0 {
		for index := len(enemySpawner.activeEnemies) - 1; index >= 0; index-- {
			enemySpawner.activeEnemies[index].Update()
			if enemySpawner.activeEnemies[index].health <= 0 {
				bank.gold += enemySpawner.activeEnemies[index].goldDropped
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
	screen.DrawImage(enemy.spritesheet, drawOps)
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
