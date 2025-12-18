package main

//Edan Hutchings
//Game Project 3

import (
	"embed"
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/paths"
	camera "github.com/tducasse/ebiten-camera"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	stage1Path = "stage1.tmx"
	stage2Path = "stage2.tmx"

	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 800

	MAP_SIZE_X = TILE_WIDTH * MAP_WIDTH
	MAP_SIZE_Y = TILE_HEIGHT * MAP_HEIGHT

	TILE_WIDTH  = 64
	TILE_HEIGHT = 64

	MAP_WIDTH  = 25
	MAP_HEIGHT = 25

	BASE_HEALTH = 100

	STARTING_GOLD = 50

	CROSSBOW_TOWER_COST = 15
)

type gameState int

const (
	gameStateStart = gameState(iota)
	gameStatePlay
	gameOver
)

type difficulty int

const (
	hard = difficulty(iota)
	easy
)

type mainGame struct {
	baseCost                int
	setDifficulty           difficulty
	state                   gameState
	ui                      *ebitenui.UI
	gameCursor              cursor
	mapGrid                 grid
	towers                  []*tower
	projManager             projectileManager
	viewX, viewY, viewSpeed int
	cameraView              *camera.Camera
	displayWorld            *ebiten.Image
	drawOps                 *ebiten.DrawImageOptions
	textOps                 *text.DrawOptions
	font                    font.Face
	enemySpawner            enemySpawn
	base                    playerBase
	bank                    goldCounter
	stageManager            stageManager
	pathMap                 *paths.Grid
}

func (game *mainGame) Update() error {
	if game.state == gameStateStart {
		game.ui.Update()
	} else if game.state == gameStatePlay {
		game.gameCursor.selectTowerType()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			buildTowerOnClick(game)
		}
		moveCamera(game)
		lockCameraInBounds(game)
		game.enemySpawner.updateEnemies(game.stageManager.currentStage.stageWaves, &game.bank, game.pathMap, &game.base)
		game.projManager.UpdateProjectiles(game.enemySpawner.activeEnemies)
		for _, tower := range game.towers {
			tower.Update(game.enemySpawner.activeEnemies, &game.projManager)
		}
		if len(game.stageManager.currentStage.stageWaves.waves) == 0 && game.stageManager.index < len(game.stageManager.stages) {
			game.stageManager.rebuildGameForStage(game)
		}
	}
	return nil
}

func lockCameraInBounds(game *mainGame) {
	if game.viewX > MAP_SIZE_X-WINDOW_WIDTH/2 {
		game.viewX = MAP_SIZE_X - WINDOW_WIDTH/2
	} else if game.viewX < WINDOW_WIDTH/2 {
		game.viewX = WINDOW_WIDTH / 2
	}
	if game.viewY > MAP_SIZE_Y-WINDOW_HEIGHT/2 {
		game.viewY = MAP_SIZE_Y - WINDOW_HEIGHT/2
	} else if game.viewY < WINDOW_HEIGHT/2 {
		game.viewY = WINDOW_HEIGHT / 2
	}
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	textFace := text.NewGoXFace(game.font)
	if game.state == gameStateStart {
		game.ui.Draw(screen)
	} else {
		game.displayWorld.DrawImage(game.stageManager.currentStage.drawableStage, game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.base.x), float64(game.base.y))
		game.displayWorld.DrawImage(game.base.sprite, game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.enemySpawner.x), float64(game.enemySpawner.y))
		game.displayWorld.DrawImage(game.enemySpawner.sprite, game.drawOps)
		game.drawOps.GeoM.Reset()
		for _, currentTower := range game.towers {
			currentTower.Draw(game.drawOps, game.displayWorld)
		}
		game.drawOps.GeoM.Reset()
		game.enemySpawner.drawEnemies(game.displayWorld, game.drawOps)
		game.projManager.DrawProjectiles(game.displayWorld, game.drawOps)
		game.textOps.GeoM.Translate(float64(game.base.x), float64(game.base.y-20))
		game.textOps.ColorScale.ScaleWithColor(colornames.White)
		text.Draw(game.displayWorld, game.base.name, textFace, game.textOps)
		game.textOps.ColorScale.Reset()
		game.textOps.GeoM.Reset()
		game.cameraView.Follow.H = game.viewY * 2
		game.cameraView.Follow.W = game.viewX * 2
		game.cameraView.Draw(game.displayWorld, screen)
		game.bank.drawCurrentGoldText(screen, game.textOps, game.font)
		game.drawOps.GeoM.Reset()
	}
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	pathMap := paths.NewGrid(25, 25, TILE_WIDTH, TILE_HEIGHT)
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	displayWorld := ebiten.NewImage(MAP_SIZE_X, MAP_SIZE_Y)
	game := mainGame{

		gameCursor:   cursor{selectedBox: nil},
		pathMap:      pathMap,
		towers:       []*tower{},
		projManager:  projectileManager{make([]projectile, 0)},
		viewX:        WINDOW_WIDTH / 2,
		viewY:        WINDOW_HEIGHT / 2,
		viewSpeed:    5,
		cameraView:   camera.Init(WINDOW_WIDTH, WINDOW_HEIGHT),
		displayWorld: displayWorld,
		stageManager: stageManager{
			stages:        getStages(),
			goToNextStage: true,
			index:         0,
		},
		drawOps: &ebiten.DrawImageOptions{},
		textOps: &text.DrawOptions{},
		font:    LoadFont("Square-Black.ttf", 30),
		bank: goldCounter{
			gold: 0, x: WINDOW_WIDTH / 2, y: 10, color: colornames.Gold,
		},
	}
	game.stageManager.buildDrawableStages()
	if game.stageManager.goToNextStage {
		game.stageManager.rebuildGameForStage(&game)
	}
	game.ui = &ebitenui.UI{Container: makeUI(&game)}
	err := ebiten.RunGame(&game)
	if err != nil {
		fmt.Println("Couldn't run game:", err)
	}
}

func moveCamera(game *mainGame) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.viewY -= game.viewSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.viewY += game.viewSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.viewX -= game.viewSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.viewX += game.viewSpeed
	}
}

func pathMaptoMapGrid(game *mainGame) {
	gridIndex := 0
	for x := 0; x < MAP_WIDTH; x += 1 {
		for y := 0; y < MAP_HEIGHT; y += 1 {
			currentGrid := game.mapGrid.gridBoxes[gridIndex]
			currentGrid.cell = game.pathMap.Get(x, y)
			if currentGrid.cell.Walkable == false {
				game.mapGrid.gridBoxes[gridIndex].canBuild = false
			}
			gridIndex += 1
		}
	}
}
