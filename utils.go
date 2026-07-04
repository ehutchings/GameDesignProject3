package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func LoadFont(fontFile string, size float64) font.Face {
	fileHandle, err := os.Open(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	fontData, err := io.ReadAll(fileHandle)
	if err != nil {
		log.Fatal(err)
	}
	ttFont, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err := opentype.NewFace(ttFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return fontFace
}

func LoadEmbeddedImage(folderName string, imageName string) *ebiten.Image {
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", folderName, imageName))
	if err != nil {
		log.Fatal("failed to load embedded image ", imageName, err)
	}
	ebitenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("Error loading tile image:", imageName, err)
	}
	return ebitenImage
}

func makeEbiteImagesFromMap(tiledMap tiled.Map) map[uint32]*ebiten.Image {
	idToImage := make(map[uint32]*ebiten.Image)
	for _, tile := range tiledMap.Tilesets[0].Tiles {
		tilePath := tile.Image.Source
		ebitenImageTile, _, err := ebitenutil.NewImageFromFile(tilePath)
		if err != nil {
			fmt.Println("Error loading tile image:", tilePath, err)
		}
		idToImage[tile.ID] = ebitenImageTile
	}
	return idToImage
}

func LoadEmbeddedWav(name string, context *audio.Context) *audio.Player {
	thunderFile, err := EmbeddedAssets.Open(path.Join("assets", "Sounds", name))
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	thunderSound, err := wav.DecodeWithoutResampling(thunderFile)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	soundPlayer, err := context.NewPlayer(thunderSound)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return soundPlayer
}

func LoadEmbeddedTileMap(name string) *tiled.Map {
	loaderOption := tiled.WithFileSystem(EmbeddedAssets)
	tileMap, err := tiled.LoadFile(path.Join("assets", "TileMaps", name), loaderOption)
	if err != nil {
		fmt.Println("Error loading tile map: ", name, err)
	}
	return tileMap
}
