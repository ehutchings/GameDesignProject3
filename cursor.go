package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type cursor struct {
	selectedBox *gridBox
	x, y        int
}

func buildTowerOnClick(game *mainGame) {
	cursorX, cursorY := ebiten.CursorPosition()
	game.gameCursor.x, game.gameCursor.y = cursorX+game.viewX-WINDOW_WIDTH/2, cursorY+game.viewY-WINDOW_HEIGHT/2
	game.mapGrid.getGridBoxAtCursor(&game.gameCursor)
	selectedGrid := game.gameCursor.selectedBox
	if selectedGrid != nil && selectedGrid.canBuild == true {
		selectedGrid.cell.Walkable = false
		if canEnemyPath(game) && game.bank.gold >= CROSSBOW_TOWER_COST {
			selectedGrid.tower = newCrossbowTower(selectedGrid.x, selectedGrid.y)
			selectedGrid.canBuild = false
			game.towers = append(game.towers, selectedGrid.tower)
			game.bank.gold -= CROSSBOW_TOWER_COST
			redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
		} else {
			selectedGrid.cell.Walkable = true
		}
	}
}
