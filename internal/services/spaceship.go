package services

import "github.com/gdamore/tcell/v2"

func GenerateShip(screenSvc *ScreenService, objects chan ScreenObject) {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		objects,
		false,
		true,
		width / 2,
		height - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
	}
	spaceship := Spaceship{baseObject, screenSvc}
	go spaceship.Move()
}

type Spaceship struct {
	BaseObject
	screenSvc *ScreenService
}

func (this *Spaceship) Move() {
	for {
		if this.IsBlocked {
			continue
		}
		switch event := this.screenSvc.GetControlEvent(); event {
		case GoLeft:
			this.X--
		case GoRight:
			this.X++
		}
		this.IsBlocked = true
		this.Objects <- this
	}
}
