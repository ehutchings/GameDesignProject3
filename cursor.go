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
	if cursorY < WINDOW_HEIGHT-BOTTOM_BAR_HEIGHT {
		game.mapGrid.getGridBoxAtCursor(&game.gameCursor)
		selectedGrid := game.gameCursor.selectedBox
		if selectedGrid != nil && selectedGrid.canBuild == true {
			selectedGrid.cell.Walkable = false
			if game.gameCursor.selectedTower == snowflake {
				buySnowflake(game, selectedGrid)
			} else if game.gameCursor.selectedTower == infernalEye {
				buyInfernalEye(game, selectedGrid)
			} else if game.gameCursor.selectedTower == crossbow {
				buyCrossbow(game, selectedGrid)
			} else if game.gameCursor.selectedTower == voidLauncher {
				buyVoidLauncher(game, selectedGrid)
			}
		}
	}
}

func buyCrossbow(game *mainGame, selectedGrid *gridBox) {
	if canEnemyPath(game) && game.bank.gold >= game.baseCost {
		newTower := newCrossbowTower(selectedGrid.x, selectedGrid.y, int(game.setDifficulty+1))
		selectedGrid.tower = &newTower
		selectedGrid.canBuild = false
		game.towers = append(game.towers, selectedGrid.tower)
		game.bank.gold -= game.baseCost
		redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
	} else {
		selectedGrid.cell.Walkable = true
	}
}

func buySnowflake(game *mainGame, selectedGrid *gridBox) {
	if canEnemyPath(game) && game.bank.gold >= game.baseCost*2 {
		newTower := newSnowflakeTower(selectedGrid.x, selectedGrid.y, int(game.setDifficulty+1))
		selectedGrid.tower = &newTower
		selectedGrid.canBuild = false
		game.towers = append(game.towers, selectedGrid.tower)
		game.bank.gold -= game.baseCost * 2
		redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
	} else {
		selectedGrid.cell.Walkable = true
	}
}

func buyInfernalEye(game *mainGame, selectedGrid *gridBox) {
	if canEnemyPath(game) && game.bank.gold >= int(float64(game.baseCost)*1.5) {
		newTower := newInfernalEyeTower(selectedGrid.x, selectedGrid.y, int(game.setDifficulty+1))
		selectedGrid.tower = &newTower
		selectedGrid.canBuild = false
		game.towers = append(game.towers, selectedGrid.tower)
		game.bank.gold -= int(float64(game.baseCost) * 1.5)
		redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
	} else {
		selectedGrid.cell.Walkable = true
	}
}

func buyVoidLauncher(game *mainGame, selectedGrid *gridBox) {
	if canEnemyPath(game) && game.bank.gold >= game.baseCost*3 {
		newTower := newVoidLauncherTower(selectedGrid.x, selectedGrid.y, int(game.setDifficulty+1))
		selectedGrid.tower = &newTower
		selectedGrid.canBuild = false
		game.towers = append(game.towers, selectedGrid.tower)
		game.bank.gold -= game.baseCost * 3
		redrawEnemyPaths(game, game.enemySpawner.activeEnemies)
	} else {
		selectedGrid.cell.Walkable = true
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
		cursor.selectedTower = snowflake
	}
}
