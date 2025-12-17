package main

import "github.com/hajimehoshi/ebiten/v2"

type cursor struct {
	selectedBox *gridBox
	x, y        int
}

func buildTowerOnClick(game *mainGame) {
	cursorX, cursorY := ebiten.CursorPosition()
	game.gameCursor.x, game.gameCursor.y = cursorX+game.viewX-WINDOW_WIDTH/2, cursorY+game.viewY-WINDOW_HEIGHT/2
	game.mapGrid.getGridBoxAtCursor(&game.gameCursor)
	selectedGrid := game.gameCursor.selectedBox
	if selectedGrid != nil && selectedGrid.tower == nil {
		selectedGrid.cell.Walkable = false
		if canEnemyPath(game) {
			selectedGrid.tower = newCrossbowTower(selectedGrid.x, selectedGrid.y)
			game.boxesWithTowers = append(game.boxesWithTowers, selectedGrid)
			redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
		} else {
			selectedGrid.cell.Walkable = true
		}
	}
}
