package services

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	exhaustGas1 = `( )
 )
( )`
	exhaustGas2 = ` ) 
( )
 ( `
 	exhaustGasTimeout = 100
)

func GenerateExhaustGas(spaceship *Spaceship, ch chan<- ScreenObject) {
	baseObject := BaseObject{
		false,
		true,
		spaceship.X + shipWidth / 2 - 1,
		spaceship.Y + shipHeight - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0,
		exhaustGas1,
		make(chan struct{}),
		make(chan struct{}),
	}
	exhaustGas := ExhaustGas{baseObject, ch, spaceship}
	go exhaustGas.Run()
}

type ExhaustGas struct {
	BaseObject
	exhaustGasChannel chan<- ScreenObject
	spaceship *Spaceship
}

func (exhaustGas *ExhaustGas) Run() {
	views := []string{exhaustGas1, exhaustGas2}
	ticker := time.NewTicker(exhaustGasTimeout * time.Millisecond)
	defer ticker.Stop()
	i := 0
	for {
		exhaustGas.exhaustGasChannel <- exhaustGas
		select {
		case <-ticker.C:
			exhaustGas.View = views[i]
			i++
			if i >= len(views) {
				i = 0
			}
		default:
		}
		if exhaustGas.spaceship.collided {
			exhaustGas.View = ""
		}
		exhaustGas.X = exhaustGas.spaceship.X  + shipWidth / 2 - 1
		exhaustGas.Y = exhaustGas.spaceship.Y + shipHeight - 1
		<-exhaustGas.UnblockCh
		if !exhaustGas.spaceship.Active {
			return
		}
	}
}