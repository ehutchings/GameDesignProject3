package main

import (
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
)

type grid struct {
	gridBoxes []*gridBox
}

type gridBox struct {
	x, y          int
	width, height int
	tower         *tower
	canBuild      bool
	cell          *paths.Cell
	tileType      string
	collider      *resolv.ConvexPolygon
}

func createGrid() *grid {
	xIndex, yIndex := 0, 0
	mapGrid := &grid{
		make([]*gridBox, 0),
	}
	for yIndex < MAP_SIZE_Y {
		newGridBox := gridBox{
			x:      xIndex,
			y:      yIndex,
			width:  TILE_WIDTH,
			height: TILE_HEIGHT,
			collider: resolv.NewRectangle(float64(xIndex+TILE_WIDTH/2), float64(yIndex+TILE_HEIGHT/2),
				TILE_WIDTH, TILE_HEIGHT),
		}
		mapGrid.gridBoxes = append(mapGrid.gridBoxes, &newGridBox)
		xIndex += TILE_WIDTH
		if xIndex >= MAP_SIZE_X {
			xIndex = 0
			yIndex += TILE_HEIGHT
		}
	}
	return mapGrid
}

func (mapGrid grid) getGridBoxAtCursor(cursor *cursor) {
	cursorCollider := resolv.NewRectangle(float64(cursor.x), float64(cursor.y), 1, 1)
	for _, box := range mapGrid.gridBoxes {
		if cursorCollider.IsContainedBy(box.collider) {
			cursor.selectedBox = box
		}
	}
}
