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
) ([]ScreenObject, chan int) {
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
		false,
	}
	go runMeteoriteCounter(ctx, &menu, template, levelName, winGoal)
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
		false,
	}
	lifeChannel := make(chan int, initialLifeNumber)
	lifeCounter := LifeCounter{baseObject, initialLifeNumber, lifeChannel}
	lifeCounter.UpdateCounterView(initialLifeNumber)
	go lifeCounter.Run(ctx, menuChan)
	return []ScreenObject{&menu, &lifeCounter}, lifeChannel
}

func runMeteoriteCounter(
	ctx context.Context,
	menu *BaseObject,
	template string,
	levelName string,
	winGoal int,
) {
	for {
		meteoriteMutx.RLock()
		menu.View = fmt.Sprintf(template, levelName, destroyedMeteorites, winGoal)
		meteoriteMutx.RUnlock()

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
		case <-ctx.Done():
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
