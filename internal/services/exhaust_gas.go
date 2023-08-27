package services

import (
	"context"
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

func GenerateExhaustGas(
	ctx context.Context,
	spaceship *Spaceship,
	ch chan<- ScreenObject,
) {
	baseObject := BaseObject{
		false,
		true,
		spaceship.X + shipWidth/2 - 1,
		spaceship.Y + shipHeight - 1,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0,
		exhaustGas1,
		make(chan struct{}),
		make(chan struct{}),
		false,
	}
	exhaustGas := ExhaustGas{baseObject, ch, spaceship}
	go exhaustGas.Run(ctx)
}

type ExhaustGas struct {
	BaseObject
	exhaustGasChannel chan<- ScreenObject
	spaceship         *Spaceship
}

func (exhaustGas *ExhaustGas) Run(ctx context.Context) {
	views := []string{exhaustGas1, exhaustGas2}
	ticker := time.NewTicker(exhaustGasTimeout * time.Millisecond)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case exhaustGas.exhaustGasChannel <- exhaustGas:
		case <-ctx.Done():
			return
		}

		select {
		case <-ticker.C:
			exhaustGas.View = views[i]
			i++
			if i >= len(views) {
				i = 0
			}
		default:
		}

		if !exhaustGas.spaceship.Vulnerable {
			exhaustGas.View = ""
		}
		exhaustGas.X = exhaustGas.spaceship.X + shipWidth/2 - 1
		exhaustGas.Y = exhaustGas.spaceship.Y + shipHeight - 1

		select {
		case <-exhaustGas.UnblockCh:
		case <-ctx.Done():
			return
		}
		if !exhaustGas.spaceship.Active {
			exhaustGas.Active = false
			return
		}
	}
}
