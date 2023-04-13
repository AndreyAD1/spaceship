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

func (this *BaseObject) Deactivate() {
	this.Active = false
}

func (this *BaseObject) Unblock() {
	this.IsBlocked = false
	this.IsDrawn = false
}

func (this *BaseObject) GetCornerCoordinates() (int, int) {
	return int(math.Round(this.X)), int(math.Round(this.Y))
}

func (this *BaseObject) GetViewCoordinates() [][]int {
	initialX, y := int(math.Round(this.X)), int(math.Round(this.Y))
	view := this.GetView()
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

func (this *BaseObject) GetStyle() tcell.Style {
	return this.Style
}

func (this *BaseObject) GetView() string {
	return this.View
}

func (this *BaseObject) MarkDrawn() {
	this.IsDrawn = true
}

func (this *BaseObject) GetDrawStatus() bool {
	return this.IsDrawn
}

func (this *BaseObject) Collide(objects []ScreenObject) {
	this.Deactivate()
}