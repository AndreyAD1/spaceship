package services

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type ScreenObject interface {
	GetCoordinates() (int, int)
	GetStyle() tcell.Style
	Move()
	Unblock()
	Deactivate()
}

type BaseObject struct {
	Objects   chan<- ScreenObject
	IsBlocked bool
	Active    bool
	X         float64
	Y         float64
	Style     tcell.Style
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
