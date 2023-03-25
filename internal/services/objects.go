package services

import "github.com/gdamore/tcell/v2"

type ScreenObject struct {
	Events    chan<- *ScreenObject
	Block     chan struct{}
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
		this.Y++
		this.Events <- this
		<-this.Block
	}
}
