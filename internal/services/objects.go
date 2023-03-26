package services

import "github.com/gdamore/tcell/v2"

type ScreenObject interface {
	GetCoordinates() (int, int)
	GetStyle() tcell.Style
	Move()
	Unblock()
	Deactivate()
}

type BaseObject struct {
	Events    chan<- ScreenObject
	IsBlocked bool
	Active    bool
	X         int
	Y         int
	Style     tcell.Style
}

func (this *BaseObject) Deactivate() {
	this.Active = false
}

func (this *BaseObject) Unblock() {
	this.IsBlocked = false
}

func (this *BaseObject) GetCoordinates() (int, int) {
	return this.X, this.Y
}

func (this *BaseObject) GetStyle() tcell.Style {
	return this.Style
}
