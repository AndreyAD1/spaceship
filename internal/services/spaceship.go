package services

import "github.com/gdamore/tcell/v2"

func GenerateShip(screenSvc *ScreenService, objects chan ScreenObject) {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		objects,
		screenSvc,
		false,
		true,
		float64(width) / 2,
		float64(height) - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.5,
	}
	spaceship := Spaceship{baseObject}
	go spaceship.Move()
}

type Spaceship struct {
	BaseObject
}

func (this *Spaceship) Move() {
	for {
		if this.IsBlocked {
			continue
		}
		switch event := this.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			this.X--
		case GoRight:
			this.X++
		}
		this.IsBlocked = true
		this.Objects <- this
	}
}
