package main

//Edan Hutchings
//Game Project 3

import (
	"embed"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	camera "github.com/tducasse/ebiten-camera"
	"golang.org/x/image/font"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	stage1Path = "stage1.tmx"

	WINDOW_WIDTH  = 500
	WINDOW_HEIGHT = 500

	MAP_SIZE_X = 64 * 20
	MAP_SIZE_Y = 64 * 20
)

type mainGame struct {
	viewX, viewY, viewSpeed int
	cameraView              *camera.Camera
	displayWorld            *ebiten.Image
	stage1Map               *tiled.Map
	drawableStage1          *ebiten.Image
	stage1TileHash          map[uint32]*ebiten.Image
	drawOps                 *ebiten.DrawImageOptions
	font                    font.Face
}

func (game *mainGame) Update() error {
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
	fmt.Println(game.viewX, game.viewY)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.displayWorld.DrawImage(game.drawableStage1, game.drawOps)
	game.drawOps.GeoM.Reset()
	game.cameraView.Follow.H = game.viewY * 2
	game.cameraView.Follow.W = game.viewX * 2
	game.lockCameraInBounds()
	game.cameraView.Draw(game.displayWorld, screen)
}

func (game *mainGame) lockCameraInBounds() {
	if game.cameraView.Follow.H < WINDOW_HEIGHT {
		game.cameraView.Follow.H = WINDOW_HEIGHT
	}
	if game.cameraView.Follow.H > MAP_SIZE_Y*2-WINDOW_HEIGHT {
		game.cameraView.Follow.H = MAP_SIZE_Y*2 - WINDOW_HEIGHT
	}
	if game.cameraView.Follow.W < WINDOW_WIDTH {
		game.cameraView.Follow.W = WINDOW_WIDTH
	}
	if game.cameraView.Follow.W > MAP_SIZE_X*2-WINDOW_WIDTH {
		game.cameraView.Follow.W = MAP_SIZE_X*2 - WINDOW_WIDTH
	}
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	stage1, err := tiled.LoadFile(stage1Path)
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	stage1Image := makeEbiteImagesFromMap(*stage1)
	displayWorld := ebiten.NewImage(MAP_SIZE_X, MAP_SIZE_Y)
	game := mainGame{
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
	buildDrawableStage(&game)
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
			ebitenTileToDraw := game.stage1TileHash[tileToDraw.ID]
			screen.DrawImage(ebitenTileToDraw, &drawOptions)
		}
	}
}
