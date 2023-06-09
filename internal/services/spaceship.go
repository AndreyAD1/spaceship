package services

import (
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
	verticalAcceleration    = 0.3
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
	objects chan ScreenObject,
	screenSvc *ScreenService,
	gameover chan *BaseObject,
	lifeChannel chan<- int,
	invulnerableChannel chan ScreenObject,
	winGoal int,
) Spaceship {
	width, height := screenSvc.screen.Size()
	baseObject := BaseObject{
		false,
		true,
		float64(width) / 2,
		float64(height) - shipHeight,
		tcell.StyleDefault.Background(tcell.ColorReset),
		maxSpeed,
		SpaceshipView,
		make(chan struct{}),
		make(chan struct{}),
	}
	spaceship := Spaceship{
		baseObject,
		objects,
		screenSvc,
		gameover,
		0,
		0,
		3,
		lifeChannel,
		false,
		invulnerableChannel,
	}
	go spaceship.Move(winGoal)
	return spaceship
}

type Spaceship struct {
	BaseObject
	Objects             chan<- ScreenObject
	ScreenSvc           *ScreenService
	gameover            chan *BaseObject
	Vx                  float64
	Vy                  float64
	lifes               int
	lifeChannel         chan<- int
	collided            bool
	invulnerableChannel chan<- ScreenObject
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

func (spaceship *Spaceship) apply_acceleration(ax, ay float64) {
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

func (spaceship *Spaceship) Move(winGoal int) {
	for {
		switch event := spaceship.ScreenSvc.GetControlEvent(); event {
		case GoLeft:
			spaceship.apply_acceleration(-horizontalAcceleration, 0)
		case GoRight:
			spaceship.apply_acceleration(horizontalAcceleration, 0)
		case GoUp:
			spaceship.apply_acceleration(0, -verticalAcceleration)
		case GoDown:
			spaceship.apply_acceleration(0, verticalAcceleration)
		case Shoot:
			go Shot(
				spaceship.ScreenSvc,
				spaceship.Objects,
				spaceship.X+2,
				spaceship.Y-1,
			)
		case NoEvent:
			spaceship.apply_acceleration(0, 0)
		}
		if destroyedMeteorites >= winGoal {
			spaceship.Deactivate()
			go DrawWin(spaceship.gameover, spaceship.ScreenSvc)
			return
		}

		if spaceship.collided {
			spaceship.invulnerableChannel <- spaceship
		} else {
			spaceship.Objects <- spaceship
		}

		select {
		case <-spaceship.Cancel:
			return
		case <-spaceship.UnblockCh:
		}
	}
}

func (spaceship *Spaceship) Collide(objects []ScreenObject) {
	if spaceship.collided {
		return
	}
	spaceship.collided = true
	spaceship.lifes--
	spaceship.lifeChannel <- spaceship.lifes
	if spaceship.lifes <= 0 {
		spaceship.Deactivate()
		go DrawGameOver(spaceship.gameover, spaceship.ScreenSvc)
		return
	}
	go spaceship.Blink()
}

func (spaceship *Spaceship) Blink() {
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
			spaceship.collided = false
			return
		}
	}
}
