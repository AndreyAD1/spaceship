package services

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func GenerateMenu(
	ctx context.Context,
	menuChan chan ScreenObject,
	levelName string,
	initialLifeNumber,
	winGoal int,
) chan int {
	go runMeteoriteCounter(ctx, menuChan, levelName, winGoal)
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	baseObject := BaseObject{
		false,
		true,
		3,
		3,
		style,
		0,
		"♥",
		make(chan (struct{})),
		make(chan (struct{})),
	}
	lifeChannel := make(chan int, initialLifeNumber)
	lifeCounter := LifeCounter{baseObject, initialLifeNumber, lifeChannel}
	lifeCounter.UpdateCounterView(initialLifeNumber)
	go lifeCounter.Run(ctx, menuChan)
	return lifeChannel
}

func runMeteoriteCounter(ctx context.Context, menuChan chan ScreenObject, levelName string, winGoal int) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	template := "%v\nDestroyed Meteorites: %v/%v"
	menu := BaseObject{
		false,
		true,
		3,
		1,
		style,
		0,
		fmt.Sprintf(template, levelName, destroyedMeteorites, winGoal),
		make(chan (struct{})),
		make(chan (struct{})),
	}
	for {
		select {
		case menuChan <- &menu:
		case <-ctx.Done():
			return
		}

		menu.View = fmt.Sprintf(template, levelName, destroyedMeteorites, winGoal)

		select {
		case <-menu.UnblockCh:
		case <-ctx.Done():
			return
		}
	}
}

type LifeCounter struct {
	BaseObject
	lifeNumber  int
	lifeChannel <-chan int
}

func (counter *LifeCounter) Run(ctx context.Context, menuChannel chan<- ScreenObject) {
	for {
		select {
		case menuChannel <- counter:
		case <-ctx.Done():
			return
		}

		select {
		case lifeNumber := <-counter.lifeChannel:
			counter.UpdateCounterView(lifeNumber)
		default:
		}

		select {
		case <-counter.UnblockCh:
		case <-ctx.Done():
			return
		}
	}
}

func (counter *LifeCounter) UpdateCounterView(lifeNumber int) {
	counter.lifeNumber = lifeNumber
	newView := "Lifes: "
	for i := 0; i < counter.lifeNumber; i++ {
		newView += "♥"
		if i < counter.lifeNumber-1 {
			newView += " "
		}
	}
	counter.View = newView
}
