package services

import "github.com/gdamore/tcell/v2"

const SpaceshipView = `  . 
 .'.
 |o|
.'o'.
|.-.|
'   '`

func GenerateShip(
	objects chan ScreenObject,
	screenSvc *ScreenService,
	gameover chan *BaseObject,
) Spaceship {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		false,
		false,
		true,
		float64(width) / 2,
		float64(height) - 6,
		tcell.StyleDefault.Background(tcell.ColorReset),
		1,
		SpaceshipView,
		make(chan struct{}),
		make(chan struct{}),
	}
	spaceship := Spaceship{baseObject, objects, screenSvc, gameover}
	go spaceship.Move()
	return spaceship
}

type Spaceship struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
	gameover  chan *BaseObject
}

func (spaceship *Spaceship) Move() {
	for {
		newX := spaceship.X
		switch event := spaceship.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			newX = spaceship.X - spaceship.Speed
		case GoRight:
			newX = spaceship.X + spaceship.Speed
		case Shoot:
			go Shot(spaceship.ScreenSvc, spaceship.Objects, spaceship.X+2, spaceship.Y-1)
		}
		leftBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(newX, spaceship.Y)
		rightBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(newX+3, spaceship.Y)
		if leftBoundaryIsValid && rightBoundaryIsValid {
			spaceship.X = newX
		}
		spaceship.Objects <- spaceship
		
		select {
		case <-spaceship.Cancel:
			go DrawGameOver(spaceship.gameover, spaceship.ScreenSvc)
			return
		case <-spaceship.UnblockCh:
		}
	}
}
