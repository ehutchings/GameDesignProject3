package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
)

type enemySpawn struct {
	sprite  *ebiten.Image
	x, y    int
	enemies []*enemy
}

type enemy struct {
	spritesheet            *ebiten.Image
	x, y                   int
	speed                  int
	xDirection, yDirection int
	path                   *paths.Path
	collider               *resolv.ConvexPolygon
}

func newEnemySpawn(x, y int) *enemySpawn {
	image := LoadEmbeddedImage("", "enemySpawn.png")
	return &enemySpawn{
		sprite:  image,
		x:       x,
		y:       y,
		enemies: make([]*enemy, 0),
	}
}

func newEnemy(x, y, speed int) *enemy {
	image := LoadEmbeddedImage("", "enemy.png")
	return &enemy{
		spritesheet: image,
		x:           x,
		y:           y,
		speed:       speed,
		xDirection:  0,
		yDirection:  0,
		collider:    resolv.NewRectangle(float64(x-TILE_WIDTH/2), float64(y-TILE_WIDTH/2), TILE_WIDTH, TILE_HEIGHT),
	}
}

func newEnemyPath(game *mainGame, enemy *enemy) {
	startingCell := game.pathMap.Get(game.enemySpawner.x/TILE_WIDTH, game.enemySpawner.y/TILE_HEIGHT)
	endingCell := game.pathMap.Get(game.base.x/TILE_WIDTH, game.base.y/TILE_HEIGHT)
	enemy.path = game.pathMap.GetPathFromCells(startingCell, endingCell, false, false)
}

func (enemy *enemy) Update() {
	if enemy.path != nil {
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
		enemy.y += enemy.yDirection * enemy.speed
		enemy.collider.SetX(float64(enemy.x - TILE_WIDTH/2))
		enemy.collider.SetY(float64(enemy.y - TILE_HEIGHT/2))
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
