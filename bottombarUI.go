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
			widget.RowLayoutOpts.Spacing(WINDOW_WIDTH/4-32),
		)), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(WINDOW_WIDTH, BOTTOM_BAR_HEIGHT)))
	backgroundImage := ebiten.NewImage(WINDOW_WIDTH, BOTTOM_BAR_HEIGHT)
	backgroundImage.Fill(colornames.Gray)
	buttonsContainer.SetLocation(image.Rect(0, WINDOW_HEIGHT-BOTTOM_BAR_HEIGHT, WINDOW_WIDTH, WINDOW_HEIGHT))
	buttonsContainer.SetBackgroundImage(imageUI.NewNineSliceSimple(backgroundImage, 10, 100))
	makeBottomBarButtons(buttonsContainer, game)
	rootContainer.AddChild(buttonsContainer)
	return rootContainer
}

func makeBottomBarButtons(container *widget.Container, game *mainGame) {
	crossbowButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Crossbow", "CrossbowDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Crossbow", "CrossbowSelected.png")),
		}),
		widget.CheckboxOpts.WidgetOpts(widget.WidgetOpts.CustomData("crossbow")),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			state := args.State
			if state == 1 {
				game.gameCursor.selectedTower = crossbow
				for _, button := range game.bottomBarButtons {
					if button.GetWidget().CustomData != "crossbow" {
						button.SetState(0)
					}
				}
			} else if game.gameCursor.selectedTower == crossbow {
				game.gameCursor.selectedTower = empty
			}
		}))
	container.AddChild(crossbowButton)
	game.bottomBarButtons = append(game.bottomBarButtons, crossbowButton)

	infernalEyeButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/InfernalEye", "InfernalEyeDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/InfernalEye", "InfernalEyeSelected.png")),
		}),
		widget.CheckboxOpts.WidgetOpts(widget.WidgetOpts.CustomData("infernalEye")),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			state := args.State
			if state == 1 {
				game.gameCursor.selectedTower = infernalEye
				for _, button := range game.bottomBarButtons {
					if button.GetWidget().CustomData != "infernalEye" {
						button.SetState(0)
					}
				}
			} else if game.gameCursor.selectedTower == infernalEye {
				game.gameCursor.selectedTower = empty
			}
		}))
	container.AddChild(infernalEyeButton)
	game.bottomBarButtons = append(game.bottomBarButtons, infernalEyeButton)

	snowflakeButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Snowflake", "SnowflakeDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/Snowflake", "SnowflakeSelected.png")),
		}),
		widget.CheckboxOpts.WidgetOpts(widget.WidgetOpts.CustomData("snowflake")),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			state := args.State
			if state == 1 {
				game.gameCursor.selectedTower = snowflake
				for _, button := range game.bottomBarButtons {
					if button.GetWidget().CustomData != "snowflake" {
						button.SetState(0)
					}
				}
			} else if game.gameCursor.selectedTower == snowflake {
				game.gameCursor.selectedTower = empty
			}
		}))
	container.AddChild(snowflakeButton)
	game.bottomBarButtons = append(game.bottomBarButtons, snowflakeButton)

	voidLauncherButton := widget.NewCheckbox(
		widget.CheckboxOpts.Image(&widget.CheckboxImage{
			Unchecked: imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/VoidLauncher", "VoidLauncherDeselected.png")),
			Checked:   imageUI.NewFixedNineSlice(LoadEmbeddedImage("TowerSelectGUI/VoidLauncher", "VoidLauncherSelected.png")),
		}),
		widget.CheckboxOpts.WidgetOpts(widget.WidgetOpts.CustomData("voidLauncher")),
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			state := args.State
			if state == 1 {
				game.gameCursor.selectedTower = voidLauncher
				for _, button := range game.bottomBarButtons {
					if button.GetWidget().CustomData != "voidLauncher" {
						button.SetState(0)
					}
				}
			} else if game.gameCursor.selectedTower == voidLauncher {
				game.gameCursor.selectedTower = empty
			}
		}))
	container.AddChild(voidLauncherButton)
	game.bottomBarButtons = append(game.bottomBarButtons, voidLauncherButton)
}
