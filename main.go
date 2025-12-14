package main

//Edan Hutchings
//Game Project 3

import (
	"embed"
	"fmt"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/paths"
	camera "github.com/tducasse/ebiten-camera"
	"golang.org/x/image/font"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	stage1Path = "stage1.tmx"

	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 800

	MAP_SIZE_X = TILE_WIDTH * MAP_WIDTH
	MAP_SIZE_Y = TILE_HEIGHT * MAP_HEIGHT

	TILE_WIDTH  = 64
	TILE_HEIGHT = 64

	MAP_WIDTH  = 25
	MAP_HEIGHT = 25
)

type gameState int

const (
	gameStateStart = gameState(iota)
	gameStatePlay
)

type mainGame struct {
	name                    string
	baseCost                int
	state                   gameState
	ui                      *ebitenui.UI
	gameCursor              cursor
	mapGrid                 *grid
	viewX, viewY, viewSpeed int
	cameraView              *camera.Camera
	displayWorld            *ebiten.Image
	stage1Map               *tiled.Map
	drawableStage1          *ebiten.Image
	stage1TileHash          map[uint32]*ebiten.Image
	drawOps                 *ebiten.DrawImageOptions
	font                    font.Face
	pathMap                 *paths.Grid
}

type enemy struct {
	spritesheet            *ebiten.Image
	x, y                   int
	xDirection, yDirection int
	path                   *paths.Path
}

type tower struct {
	spritesheet *ebiten.Image
	x, y        int
	baseDamage  int
	rangeRadius int
	baseCostMod float64
}

func (game *mainGame) Update() error {
	if game.state == gameStateStart {
		game.ui.Update()
	} else if game.state == gameStatePlay {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			cursorX, cursorY := ebiten.CursorPosition()
			game.gameCursor.x, game.gameCursor.y = cursorX+game.viewX-WINDOW_WIDTH/2, cursorY+game.viewY-WINDOW_HEIGHT/2
			game.mapGrid.getGridBoxAtCursor(&game.gameCursor)
		}
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
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	if game.state == gameStateStart {
		game.ui.Draw(screen)
	} else {
		game.displayWorld.DrawImage(game.drawableStage1, game.drawOps)
		game.drawOps.GeoM.Reset()
		game.cameraView.Follow.H = game.viewY * 2
		game.cameraView.Follow.W = game.viewX * 2
		game.cameraView.Draw(game.displayWorld, screen)
	}
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	pathMap := paths.NewGrid(25, 25, TILE_WIDTH, TILE_HEIGHT)
	stage1, err := tiled.LoadFile(stage1Path)
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	stage1Image := makeEbiteImagesFromMap(*stage1)
	displayWorld := ebiten.NewImage(MAP_SIZE_X, MAP_SIZE_Y)
	game := mainGame{
		pathMap:        pathMap,
		gameCursor:     cursor{selectedBox: nil},
		mapGrid:        createGrid(),
		viewX:          WINDOW_WIDTH / 2,
		viewY:          WINDOW_HEIGHT / 2,
		viewSpeed:      5,
		cameraView:     camera.Init(WINDOW_WIDTH, WINDOW_HEIGHT),
		displayWorld:   displayWorld,
		stage1Map:      stage1,
		drawableStage1: ebiten.NewImage(stage1.Width*stage1.TileWidth, stage1.Height*stage1.TileHeight),
		stage1TileHash: stage1Image,
		drawOps:        &ebiten.DrawImageOptions{},
		font:           LoadFont("Square-Black.ttf", 30),
	}
	game.ui = &ebitenui.UI{Container: makeUI(&game)}
	buildDrawableStage(&game)
	buildPathMap(&game)
	pathMaptoMapGrid(&game)
	err = ebiten.RunGame(&game)
	if err != nil {
		fmt.Println("Couldn't run game:", err)
	}
}

func buildDrawableStage(game *mainGame) {
	screen := game.drawableStage1
	drawOptions := ebiten.DrawImageOptions{}
	for tileY := 0; tileY < game.stage1Map.Height; tileY += 1 {
		for tileX := 0; tileX < game.stage1Map.Width; tileX += 1 {
			drawOptions.GeoM.Reset()
			TileXpos := float64(game.stage1Map.TileWidth * tileX)
			TileYpos := float64(game.stage1Map.TileHeight * tileY)
			drawOptions.GeoM.Translate(TileXpos, TileYpos)
			tileToDraw := game.stage1Map.Layers[0].Tiles[tileY*game.stage1Map.Width+tileX]
			if tileToDraw.ID != 0 {
				ebitenTileToDraw := game.stage1TileHash[tileToDraw.ID]
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
		}
	}
}

func buildPathMap(game *mainGame) {
	for tileY := 0; tileY < game.stage1Map.Height; tileY += 1 {
		for tileX := 0; tileX < game.stage1Map.Width; tileX += 1 {
			currentTile := game.stage1Map.Layers[0].Tiles[tileY*game.stage1Map.Width+tileX]
			if currentTile.ID == 3 {
				game.pathMap.Get(tileX, tileY).Walkable = false
			}
		}
	}
}

func pathMaptoMapGrid(game *mainGame) {
	gridIndex := 0
	for x := 0; x < MAP_WIDTH; x += 1 {
		for y := 0; y < MAP_HEIGHT; y += 1 {
			currentGrid := game.mapGrid.gridBoxes[gridIndex]
			currentGrid.cell = game.pathMap.Get(y, x) //Grid is built by columns first, so flip x and y
			gridIndex += 1
		}
	}
}
