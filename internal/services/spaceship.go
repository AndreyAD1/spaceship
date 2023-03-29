package services

import "github.com/gdamore/tcell/v2"

func GenerateShip(objects chan ScreenObject, screenSvc *ScreenService) {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		objects,
		screenSvc,
		false,
		true,
		float64(width) / 2,
		float64(height) - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
		1,
	}
	spaceship := Spaceship{baseObject}
	go spaceship.Move()
}

type Spaceship struct {
	BaseObject
}

func (this *Spaceship) Move() {
	for {
		if !this.Active {
			// go GameOver(this.Objects)
			break
		}
		if this.IsBlocked {
			continue
		}
		newX := this.X
		switch event := this.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			newX = this.X - this.Speed
		case GoRight:
			newX = this.X + this.Speed
		}
		if this.ScreenSvc.IsInsideScreen(newX, this.Y) {
			this.X = newX
		}
		this.IsBlocked = true
		this.Objects <- this
	}
}