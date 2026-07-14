package main

import (
	"image"

	imageUI "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const BOTTOM_BAR_HEIGHT = 120

func makeBottomBarUI(game *mainGame) *widget.Container {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(&widget.Insets{
				Top:    WINDOW_HEIGHT - BOTTOM_BAR_HEIGHT,
				Left:   0,
				Right:  0,
				Bottom: 0,
			}),
			widget.RowLayoutOpts.Direction(
				widget.DirectionHorizontal,
			),
			widget.RowLayoutOpts.Spacing(5),
		)))

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Direction(
				widget.DirectionHorizontal,
			),
			widget.RowLayoutOpts.Spacing(5),
		)), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(WINDOW_WIDTH, BOTTOM_BAR_HEIGHT)))
	backgroundImage := ebiten.NewImage(WINDOW_WIDTH, BOTTOM_BAR_HEIGHT)
	backgroundImage.Fill(colornames.Gray)
	buttonsContainer.SetLocation(image.Rect(0, WINDOW_HEIGHT-BOTTOM_BAR_HEIGHT, WINDOW_WIDTH, WINDOW_HEIGHT))
	buttonsContainer.SetBackgroundImage(imageUI.NewNineSliceSimple(backgroundImage, 10, 100))
	makeBottomBarButtons(buttonsContainer)
	rootContainer.AddChild(buttonsContainer)
	return rootContainer
}

func makeBottomBarButtons(container *widget.Container) {
	crossbowButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Crossbow", "CrossbowDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Crossbow", "CrossbowSelected.png")),
		}))
	container.AddChild(crossbowButton)
	infernalEyeButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/InfernalEye", "InfernalEyeDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/InfernalEye", "InfernalEyeSelected.png")),
		}))
	container.AddChild(infernalEyeButton)
	snowflakeButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Snowflake", "SnowflakeDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Snowflake", "SnowflakeSelected.png")),
		}))
	container.AddChild(snowflakeButton)
	voidLauncherButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/VoidLauncher", "VoidLauncherDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/VoidLauncher", "VoidLauncherSelected.png")),
		}))
	container.AddChild(voidLauncherButton)
}
