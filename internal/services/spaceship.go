package services

import "github.com/gdamore/tcell/v2"

func GenerateShip(
	objects chan ScreenObject, 
	screenSvc *ScreenService,
	gameover chan *BaseObject,
) Spaceship {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		false,
		true,
		float64(width) / 2,
		float64(height) - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
		1,
		"O",
	}
	spaceship := Spaceship{baseObject, objects, screenSvc, gameover}
	go spaceship.Move()
	return spaceship
}

type Spaceship struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
	gameover chan *BaseObject
}

func (this *Spaceship) Move() {
	for {
		if !this.Active {
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
		case Shoot:
			go Shot(this.ScreenSvc, this.Objects, this.X, this.Y - 1)
		}
		if this.ScreenSvc.IsInsideScreen(newX, this.Y) {
			this.X = newX
		}
		this.IsBlocked = true
		this.Objects <- this
	}
	go DrawGameOver(this.gameover, this.ScreenSvc)
}
