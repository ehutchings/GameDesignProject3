package main

import "github.com/solarlune/paths"

type stageManager struct {
	stages        []*stage
	currentStage  *stage
	goToNextStage bool
	index         int
}

func (manager *stageManager) buildDrawableStages() {
	for _, stage := range manager.stages {
		stage.buildDrawableStage()
	}
}

func (manager *stageManager) setStageAtIndex() {
	manager.currentStage = manager.stages[manager.index]
	manager.index++
	manager.goToNextStage = false
}

func (manager *stageManager) rebuildGameForStage(game *mainGame) {
	game.mapGrid = createGrid()
	game.pathMap = paths.NewGrid(25, 25, TILE_WIDTH, TILE_HEIGHT)
	manager.setStageAtIndex()
	game.enemySpawner = *newEnemySpawn(manager.currentStage.stageWaves.getNextWave(),
		manager.currentStage.enemySpawnX, manager.currentStage.enemySpawnY)
	game.base = *newPlayerBase(manager.currentStage.playerBaseX, manager.currentStage.playerBaseY)
	manager.currentStage.buildPathMap(game.pathMap)
	pathMaptoMapGrid(game)
	game.towers = game.towers[:0]
	game.bank.gold = STARTING_GOLD
}
