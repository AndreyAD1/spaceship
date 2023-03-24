package services

import "github.com/gdamore/tcell/v2"

type ScreenObject struct {
	Events    chan<- *ScreenObject
	Block     chan struct{}
	Eliminate chan struct{}
	X         int
	Y         int
	Style     tcell.Style
}

func (this *ScreenObject) Move() {
	for {
		select {
		case <-this.Eliminate:
			break
		default:
		}
		this.Y++
		this.Events <- this
		<-this.Block
	}
}