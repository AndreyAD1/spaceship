package services

import (
	"math"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type ScreenObject interface {
	GetCoordinates() [][]int
	GetStyle() tcell.Style
	Unblock()
	Deactivate()
	GetView() string
	GetDrawStatus() bool
	MarkDrawn()
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

func (this *BaseObject) GetCoordinates() [][]int {
	initialX, y := int(math.Round(this.X)), int(math.Round(this.Y))
	view := this.GetView()
	x := initialX
	coordinates := [][]int{}
	for _, char := range view {
		if unicode.IsControl(char) {
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
