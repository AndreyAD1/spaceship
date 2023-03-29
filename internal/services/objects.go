package services

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type ScreenObject interface {
	GetCoordinates() (int, int)
	GetStyle() tcell.Style
	Unblock()
	Deactivate()
}

type BaseObject struct {
	IsBlocked bool
	Active    bool
	X         float64
	Y         float64
	Style     tcell.Style
	Speed     float64 // Cells per iteration. Max speed = 1
}

func (this *BaseObject) Deactivate() {
	this.Active = false
}

func (this *BaseObject) Unblock() {
	this.IsBlocked = false
}

func (this *BaseObject) GetCoordinates() (int, int) {
	return int(math.Round(this.X)), int(math.Round(this.Y))
}

func (this *BaseObject) GetStyle() tcell.Style {
	return this.Style
}
