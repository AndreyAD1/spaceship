package services

import (
	"context"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	SpaceshipView = `  . 
 .'.
 |o|
.'o'.
|.-.|
'   '`
	frictionCoefficient     = 0.9
	accelerationCoefficient = 0.75
	shipHeight              = 6
	shipWidth               = 5
	maxSpeed                = 0.3
	verticalAcceleration    = 0.45
	horizontalAcceleration  = 0.8
	blinkDuration           = 2000
	blinkTimeout            = 200
)

var emptyView string = string([]rune{
	' ', ' ', 0x85, ' ', '\n',
	' ', 0x85, 0x85, 0x85, ' ', '\n',
	' ', 0x85, 0x85, 0x85, ' ', '\n',
	0x85, 0x85, 0x85, 0x85, 0x85, '\n',
	0x85, 0x85, 0x85, 0x85, 0x85, '\n',
	0x85, ' ', ' ', ' ', 0x85, '\n',
})

func GenerateShip(
	ctx context.Context,
	objects chan ScreenObject,
	screenSvc *ScreenService,
	lifeChannel chan<- int,
	initialLifeNumber int,
) Spaceship {
	width, height := screenSvc.GetScreenSize()
	baseObject := BaseObject{
		false,
		true,
		float64(width) / 2,
		float64(height) - shipHeight - 3,
		tcell.StyleDefault.Background(tcell.ColorReset),
		maxSpeed,
		SpaceshipView,
		make(chan struct{}),
		make(chan struct{}),
		true,
	}
	spaceship := Spaceship{
		baseObject,
		objects,
		screenSvc,
		0,
		0,
		initialLifeNumber,
		lifeChannel,
	}
	go spaceship.Move(ctx)
	go GenerateExhaustGas(ctx, &spaceship, objects)
	return spaceship
}

type Spaceship struct {
	BaseObject
	Objects     chan<- ScreenObject
	ScreenSvc   *ScreenService
	Vx          float64
	Vy          float64
	Lifes       int
	lifeChannel chan<- int
}

func (spaceship *Spaceship) getNewSpeed(
	velocity,
	acceleration,
	frictionCoef float64,
) float64 {
	speedRate := velocity / spaceship.MaxSpeed
	accelerationRate := math.Cos(speedRate) * accelerationCoefficient
	newSpeed := velocity + acceleration*accelerationRate
	return newSpeed * frictionCoef
}

func (spaceship *Spaceship) applyAcceleration(ax, ay float64) {
	spaceship.Vx = spaceship.getNewSpeed(spaceship.Vx, ax, frictionCoefficient)
	spaceship.Vy = spaceship.getNewSpeed(spaceship.Vy, ay, frictionCoefficient)
	newX := spaceship.X + spaceship.Vx
	newY := spaceship.Y + spaceship.Vy
	leftBoundaryIsOk := spaceship.ScreenSvc.IsInsideScreen(newX, spaceship.Y)
	rightBoundaryIsOk := spaceship.ScreenSvc.IsInsideScreen(
		newX+shipWidth-2,
		spaceship.Y,
	)
	if !leftBoundaryIsOk || !rightBoundaryIsOk {
		spaceship.Vx = 0
		newX = spaceship.X
	}
	upperBoundaryIsOk := spaceship.ScreenSvc.IsInsideScreen(spaceship.X, newY)
	lowerBoundaryIsOk := spaceship.ScreenSvc.IsInsideScreen(
		spaceship.X,
		newY+shipHeight-1,
	)
	if !upperBoundaryIsOk || !lowerBoundaryIsOk {
		spaceship.Vy = 0
		newY = spaceship.Y
	}
	spaceship.X = newX
	spaceship.Y = newY
}

func (spaceship *Spaceship) Move(ctx context.Context) {
	select {
	case spaceship.Objects <- spaceship:
	case <-ctx.Done():
		return
	}

	for {
		switch event := spaceship.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			spaceship.applyAcceleration(-horizontalAcceleration, 0)
		case GoRight:
			spaceship.applyAcceleration(horizontalAcceleration, 0)
		case GoUp:
			spaceship.applyAcceleration(0, -verticalAcceleration)
		case GoDown:
			spaceship.applyAcceleration(0, verticalAcceleration)
		case Shoot:
			if spaceship.Vulnerable {
				go Shot(
					ctx,
					spaceship.ScreenSvc,
					spaceship.Objects,
					spaceship.X+2,
					spaceship.Y-1,
				)
			}
		case NoEvent:
			spaceship.applyAcceleration(0, 0)
		}

		select {
		case <-ctx.Done():
			return
		case <-spaceship.Cancel:
			return
		case <-spaceship.UnblockCh:
		}
	}
}

func (spaceship *Spaceship) Collide(ctx context.Context, objects []ScreenObject) bool {
	if !spaceship.Vulnerable {
		return false
	}
	spaceship.Vulnerable = false
	spaceship.Lifes--
	spaceship.lifeChannel <- spaceship.Lifes
	if spaceship.Lifes > 0 {
		go spaceship.Blink(ctx)
		return true
	}
	spaceship.Deactivate()
	go Explode(ctx, spaceship.Objects, spaceship.X, spaceship.Y)
	return true
}

func (spaceship *Spaceship) Blink(ctx context.Context) {
	views := []string{emptyView, SpaceshipView}
	ticker := time.NewTicker(blinkTimeout * time.Millisecond)
	defer ticker.Stop()
	spaceship.View = emptyView
	abort := time.After(blinkDuration * time.Millisecond)
	i := 0
	for {
		select {
		case <-ticker.C:
			spaceship.View = views[i%2]
			i++
		case <-abort:
			spaceship.View = SpaceshipView
			spaceship.Vulnerable = true
			return
		case <-ctx.Done():
			return
		}
	}
}
