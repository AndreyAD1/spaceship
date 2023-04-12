package services

import "github.com/gdamore/tcell/v2"

func Shot(screenSvc *ScreenService, objects chan<- ScreenObject, x, y float64) {
	baseObject := BaseObject{
		false,
		false,
		true,
		x,
		y,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.1,
		"|",
	}
	shot := Shell{baseObject, objects, screenSvc}
	go shot.Move()
}

type Shell struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
}


func (this *Shell) Move() {
	for {
		if this.Active != true {
			break
		}
		if this.IsBlocked {
			continue
		}
		newY := this.Y - this.Speed
		if !this.ScreenSvc.IsInsideScreen(this.X, newY) {
			this.Deactivate()
			break
		} 
		this.Y = newY
		this.IsBlocked = true
		this.Objects <- this
	}
}
