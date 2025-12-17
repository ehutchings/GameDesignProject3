package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type cursor struct {
	selectedBox   *gridBox
	selectedTower towerType
	x, y          int
}

func buildTowerOnClick(game *mainGame) {
	cursorX, cursorY := ebiten.CursorPosition()
	game.gameCursor.x, game.gameCursor.y = cursorX+game.viewX-WINDOW_WIDTH/2, cursorY+game.viewY-WINDOW_HEIGHT/2
	game.mapGrid.getGridBoxAtCursor(&game.gameCursor)
	selectedGrid := game.gameCursor.selectedBox
	if selectedGrid != nil && selectedGrid.canBuild == true {
		selectedGrid.cell.Walkable = false
		if game.gameCursor.selectedTower == crossbow {
			if canEnemyPath(game) && game.bank.gold >= CROSSBOW_TOWER_COST {
				newTower := newCrossbowTower(selectedGrid.x, selectedGrid.y)
				selectedGrid.tower = &newTower
				selectedGrid.canBuild = false
				game.towers = append(game.towers, selectedGrid.tower)
				game.bank.gold -= CROSSBOW_TOWER_COST
				redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
			}
		} else if game.gameCursor.selectedTower == voidLauncher {
			if canEnemyPath(game) && game.bank.gold >= CROSSBOW_TOWER_COST {
				newTower := newVoidLauncherTower(selectedGrid.x, selectedGrid.y)
				selectedGrid.tower = &newTower
				selectedGrid.canBuild = false
				game.towers = append(game.towers, selectedGrid.tower)
				game.bank.gold -= CROSSBOW_TOWER_COST
				redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
			} else {
				selectedGrid.cell.Walkable = true
			}
		}
	}
}

func (cursor *cursor) selectTowerType() {
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		cursor.selectedTower = crossbow
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		cursor.selectedTower = voidLauncher
	}
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		cursor.selectedTower = infernalEye
	}
	if ebiten.IsKeyPressed(ebiten.KeyB) {
		//TODO
	}
}
