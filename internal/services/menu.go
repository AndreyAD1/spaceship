package services

import (
	"github.com/gdamore/tcell/v2"
)

func GenerateMenu(
	menuChan chan ScreenObject,
	screenSvc *ScreenService,
) chan int {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	baseObject := BaseObject{
		false,
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
	initialLifeNumber := 3
	lifeChannel := make(chan int)
	lifeCounter := LifeCounter{baseObject, initialLifeNumber, lifeChannel}
	lifeCounter.UpdateCounterView(initialLifeNumber)
	go lifeCounter.Run(menuChan)
	return lifeChannel
}

type LifeCounter struct {
	BaseObject
	lifeNumber  int
	lifeChannel <-chan int
}

func (counter *LifeCounter) Run(menuChannel chan<- ScreenObject) {
	for {
		select {
		case lifeNumber := <-counter.lifeChannel:
			counter.UpdateCounterView(lifeNumber)
		case menuChannel <- counter:
		}
	}
}

func (counter *LifeCounter) UpdateCounterView(lifeNumber int) {
	counter.lifeNumber = lifeNumber
	newView := ""
	for i := 0; i < counter.lifeNumber; i++ {
		newView += "♥"
		if i < counter.lifeNumber-1 {
			newView += " "
		}
	}
	counter.View = newView
}
