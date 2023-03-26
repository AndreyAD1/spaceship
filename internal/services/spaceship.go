package services

import "github.com/gdamore/tcell/v2"

func GenerateShip(screenSvc ScreenService, events chan ScreenObject) {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject {
		events,
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
	screenSvc ScreenService
}

func (this *Spaceship) Move() {
	for {
		if this.IsBlocked {
			continue
		}
		// TODO make a ship react to every key click
		switch this.screenSvc.GetScreenEvent() {
		case GoLeft:
			this.X--
		case GoRight:
			this.X++
		}
		this.IsBlocked = true
		this.Events <- this
	}
}