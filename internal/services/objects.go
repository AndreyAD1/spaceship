package services

import (
	"math"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type ScreenObject interface {
	GetCornerCoordinates() (int, int)
	GetViewCoordinates() [][]int
	GetStyle() tcell.Style
	Unblock()
	Deactivate()
	IsActive() bool
	GetView() string
	GetDrawStatus() bool
	MarkDrawn()
	Collide([]ScreenObject)
}

type BaseObject struct {
	IsBlocked bool
	IsDrawn   bool
	Active    bool
	X         float64 // a column of left upper corner
	Y         float64 // a row of left upper corner
	Style     tcell.Style
	Speed     float64 // Cells per iteration. Max speed = 1
	View      string
}

func (baseObject *BaseObject) Deactivate() {
	baseObject.Active = false
}

func (baseObject *BaseObject) IsActive() bool {
	return baseObject.Active
}

func (baseObject *BaseObject) Unblock() {
	baseObject.IsBlocked = false
	baseObject.IsDrawn = false
}

func (baseObject *BaseObject) GetCornerCoordinates() (int, int) {
	return int(math.Round(baseObject.X)), int(math.Round(baseObject.Y))
}

func (baseObject *BaseObject) GetViewCoordinates() [][]int {
	initialX, y := int(math.Round(baseObject.X)), int(math.Round(baseObject.Y))
	view := baseObject.GetView()
	x := initialX
	coordinates := [][]int{}
	for _, char := range view {
		if char == '\n' {
			y++
			x = initialX
			continue
		}
		if !unicode.IsSpace(char) {
			coordinates = append(coordinates, []int{x, y})
		}
		x++
	}
	return coordinates
}

func (baseObject *BaseObject) GetStyle() tcell.Style {
	return baseObject.Style
}

func (baseObject *BaseObject) GetView() string {
	return baseObject.View
}

func (baseObject *BaseObject) MarkDrawn() {
	baseObject.IsDrawn = true
}

func (baseObject *BaseObject) GetDrawStatus() bool {
	return baseObject.IsDrawn
}

func (baseObject *BaseObject) Collide(objects []ScreenObject) {
	baseObject.Deactivate()
}
