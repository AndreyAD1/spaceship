package services

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

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
		true,
		float64(width) / 2,
		float64(height) - 6,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.9,
		SpaceshipView,
		make(chan struct{}),
		make(chan struct{}),
	}
	spaceship := Spaceship{baseObject, objects, screenSvc, gameover, 0, 0}
	go spaceship.Move()
	return spaceship
}

type Spaceship struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
	gameover  chan *BaseObject
	Vx float64
	Vy float64
}

func (spaceship *Spaceship) apply_acceleration(ax, ay float64) {
	xSpeedRate := spaceship.Vx / spaceship.MaxSpeed
	ySpeedRate := spaceship.Vx / spaceship.MaxSpeed

	axRate := math.Cos(xSpeedRate) * 0.75
	ayRate := math.Cos(ySpeedRate) * 0.75

	spaceship.Vx += ax * axRate
	spaceship.Vy += ay * ayRate

	spaceship.Vx *= 0.9
	spaceship.Vy *= 0.95

	// change coordinates
	newX := spaceship.X + spaceship.Vx
	newY := spaceship.Y + spaceship.Vy
	leftBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(newX, spaceship.Y)
	rightBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(newX+3, spaceship.Y)
	if !leftBoundaryIsValid || !rightBoundaryIsValid {
		spaceship.Vx = 0
		newX = spaceship.X
	}
	upperBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(spaceship.X, newY)
	lowerBoundaryIsValid := spaceship.ScreenSvc.IsInsideScreen(spaceship.X, newY+5)
	if !upperBoundaryIsValid || !lowerBoundaryIsValid {
		spaceship.Vy = 0
		newY = spaceship.Y
	}
	spaceship.X = newX
	spaceship.Y = newY
}

func (spaceship *Spaceship) Move() {
	for {
		switch event := spaceship.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			spaceship.apply_acceleration(-0.8, 0)
		case GoRight:
			spaceship.apply_acceleration(0.8, 0)
		case GoUp:
			spaceship.apply_acceleration(0, -0.1)
		case GoDown:
			spaceship.apply_acceleration(0, 0.1)
		case Shoot:
			go Shot(spaceship.ScreenSvc, spaceship.Objects, spaceship.X+2, spaceship.Y-1)
		case NoEvent:
			spaceship.apply_acceleration(0, 0)
		}
		spaceship.Objects <- spaceship

		select {
		case <-spaceship.Cancel:
			spaceship.Active = false
			go DrawGameOver(spaceship.gameover, spaceship.ScreenSvc)
			return
		case <-spaceship.UnblockCh:
		}
	}
}
