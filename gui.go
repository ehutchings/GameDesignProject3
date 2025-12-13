package main

import (
	"bytes"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

func makeUI(game *mainGame) *widget.Container {
	textFont := DefaultFont(20)
	textInputImage := LoadEmbeddedImage("", "uiImage.png")
	buttonDefaultImage := LoadEmbeddedImage("", "uiButtonDefault.png")
	buttonPressedImage := LoadEmbeddedImage("", "uiButtonPressed.png")
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(100)),
			widget.RowLayoutOpts.Direction(
				widget.DirectionVertical,
			),
			widget.RowLayoutOpts.Spacing(5),
		)))

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))

	basePriceInputLabel := widget.NewLabel(
		widget.LabelOpts.LabelColor(&widget.LabelColor{
			Idle:     colornames.White,
			Disabled: colornames.Gray,
		}),
		widget.LabelOpts.LabelFace(&textFont))
	basePriceInputLabel.Label = "Enter the base cost of your towers"

	basePriceInput := widget.NewTextInput(
		widget.TextInputOpts.Face(&textFont),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          colornames.Bisque,
			Disabled:      colornames.Gray,
			Caret:         colornames.Black,
			DisabledCaret: colornames.Gray,
		}),
		widget.TextInputOpts.Image(
			&widget.TextInputImage{
				Idle:      image.NewNineSliceBorder(textInputImage, 20),
				Disabled:  image.NewNineSliceBorder(textInputImage, 20),
				Highlight: image.NewNineSliceBorder(textInputImage, 20),
			}),
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 30)),
		widget.TextInputOpts.Padding(&widget.Insets{
			Top:    5,
			Bottom: 5,
			Left:   10,
			Right:  10,
		}))

	playerNameInputLabel := widget.NewLabel(
		widget.LabelOpts.LabelColor(&widget.LabelColor{
			Idle:     colornames.White,
			Disabled: colornames.Gray,
		}),
		widget.LabelOpts.LabelFace(&textFont))
	playerNameInputLabel.Label = "Please enter your name"

	playerNameInput := widget.NewTextInput(
		widget.TextInputOpts.Face(&textFont),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          colornames.Bisque,
			Disabled:      colornames.Gray,
			Caret:         colornames.Black,
			DisabledCaret: colornames.Gray,
		}),
		widget.TextInputOpts.Image(
			&widget.TextInputImage{
				Idle:      image.NewNineSliceBorder(textInputImage, 20),
				Disabled:  image.NewNineSliceBorder(textInputImage, 20),
				Highlight: image.NewNineSliceBorder(textInputImage, 20),
			}),
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 30)),
		widget.TextInputOpts.Padding(&widget.Insets{
			Top:    5,
			Bottom: 5,
			Left:   10,
			Right:  10,
		}))

	startGame := func(args *widget.ButtonClickedEventArgs) {
		baseTowerCost, _ := strconv.ParseInt(basePriceInput.GetText(), 10, 64)
		playerName := playerNameInput.GetText()
		if baseTowerCost >= 0 {
			game.state = gameStatePlay
			game.name = playerName
			game.baseCost = int(baseTowerCost)
		}
	}

	button := widget.NewButton(
		widget.ButtonOpts.TextLabel("Start Game"),
		widget.ButtonOpts.ClickedHandler(startGame),
		widget.ButtonOpts.TextFace(&textFont),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:     colornames.White,
			Disabled: colornames.Gray,
			Hover:    colornames.Gold,
			Pressed:  colornames.Gold,
		}),
		widget.ButtonOpts.Image(
			&widget.ButtonImage{
				Idle:            image.NewNineSliceBorder(buttonDefaultImage, 32),
				Hover:           image.NewNineSliceBorder(buttonDefaultImage, 32),
				Pressed:         image.NewNineSliceBorder(buttonPressedImage, 32),
				PressedHover:    image.NewNineSliceBorder(buttonPressedImage, 32),
				Disabled:        image.NewNineSliceBorder(buttonDefaultImage, 32),
				PressedDisabled: nil,
			}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Bottom: 30,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 60)))

	rootContainer.AddChild(basePriceInputLabel)
	rootContainer.AddChild(basePriceInput)
	rootContainer.AddChild(playerNameInputLabel)
	rootContainer.AddChild(playerNameInput)
	buttonContainer.AddChild(button)
	rootContainer.AddChild(buttonContainer)
	return rootContainer
}

func DefaultFont(size float64) text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}
