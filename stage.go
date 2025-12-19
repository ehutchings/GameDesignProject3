package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/paths"
)

type stageNumber int

const (
	stage1 stageNumber = iota
	stage2
	stage3
)

func stageNumberToPath(number stageNumber) string {
	switch number {
	case stage1:
		return stage1Path
	case stage2:
		return stage2Path
	case stage3:
		return stage3Path
	default:
		return ""
	}
}

type stage struct {
	stageMap                 *tiled.Map
	drawableStage            *ebiten.Image
	stageTileHash            map[uint32]*ebiten.Image
	number                   stageNumber
	stageWaves               *stageWaves
	enemySpawnX, enemySpawnY int
	playerBaseX, playerBaseY int
}

func (stage *stage) buildDrawableStage() {
	stageMap, _ := tiled.LoadFile(stageNumberToPath(stage.number))
	stage.stageMap = stageMap
	stageImage := makeEbiteImagesFromMap(*stage.stageMap)
	stage.stageTileHash = stageImage
	stage.drawableStage = ebiten.NewImage(stage.stageMap.Width*stage.stageMap.TileWidth,
		stage.stageMap.Height*stage.stageMap.TileHeight)
	screen := stage.drawableStage
	drawOptions := ebiten.DrawImageOptions{}
	for tileY := 0; tileY < stage.stageMap.Height; tileY += 1 {
		for tileX := 0; tileX < stage.stageMap.Width; tileX += 1 {
			drawOptions.GeoM.Reset()
			TileXpos := float64(stage.stageMap.TileWidth * tileX)
			TileYpos := float64(stage.stageMap.TileHeight * tileY)
			drawOptions.GeoM.Translate(TileXpos, TileYpos)
			tileToDraw := stage.stageMap.Layers[0].Tiles[tileY*stage.stageMap.Width+tileX]
			if tileToDraw.ID != 0 {
				ebitenTileToDraw := stage.stageTileHash[tileToDraw.ID]
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
		}
	}
}

func (stage *stage) buildPathMap(pathMap *paths.Grid) {
	for tileY := 0; tileY < stage.stageMap.Height; tileY += 1 {
		for tileX := 0; tileX < stage.stageMap.Width; tileX += 1 {
			currentTile := stage.stageMap.Layers[0].Tiles[tileY*stage.stageMap.Width+tileX]
			if currentTile.ID == 3 || currentTile.ID == 4 {
				pathMap.Get(tileX, tileY).Walkable = false
			}
		}
	}
}

func getStages() []*stage {
	stage1 := stage{
		stageMap:      nil,
		drawableStage: nil,
		stageTileHash: nil,
		number:        stage1,
		enemySpawnX:   0,
		enemySpawnY:   0,
		playerBaseX:   24 * TILE_WIDTH,
		playerBaseY:   20 * TILE_HEIGHT,
	}
	stage1.stageWaves = newWavesForStage(60, 2, 5)
	stage2 := stage{
		stageMap:      nil,
		drawableStage: nil,
		stageTileHash: nil,
		number:        stage2,
		enemySpawnX:   0,
		enemySpawnY:   14 * TILE_HEIGHT,
		playerBaseX:   23 * TILE_WIDTH,
		playerBaseY:   3 * TILE_HEIGHT,
	}
	stage2.stageWaves = newWavesForStage(60, 2, 10)
	stage3 := stage{
		stageMap:      nil,
		drawableStage: nil,
		stageTileHash: nil,
		number:        stage3,
		enemySpawnX:   10 * TILE_WIDTH,
		enemySpawnY:   22 * TILE_HEIGHT,
		playerBaseX:   18 * TILE_WIDTH,
		playerBaseY:   10 * TILE_HEIGHT,
	}
	stage3.stageWaves = newWavesForStage(30, 3, 20)
	return []*stage{&stage1, &stage2, &stage3}
}
