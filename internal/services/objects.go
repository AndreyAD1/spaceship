package services

import "github.com/gdamore/tcell/v2"

type ScreenObject struct {
	Events    chan<- *ScreenObject
	IsBlocked bool
	Active    bool
	X         int
	Y         int
	Style     tcell.Style
}

func (this *ScreenObject) Move() {
	for {
		if this.Active != true {
			break
		}
		if this.IsBlocked {
			continue
		}
		this.Y++
		this.IsBlocked = true
		this.Events <- this
	}
}
