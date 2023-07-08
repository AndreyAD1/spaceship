package services

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func GenerateMenu(menuChan chan ScreenObject, winGoal int) chan int {
	go runGoalMenu(menuChan, winGoal)
	go runMeteoriteCounter(menuChan)
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
	initialLifeNumber := 3
	lifeChannel := make(chan int)
	lifeCounter := LifeCounter{baseObject, initialLifeNumber, lifeChannel}
	lifeCounter.UpdateCounterView(initialLifeNumber)
	go lifeCounter.Run(menuChan)
	return lifeChannel
}

func runGoalMenu(menuChan chan ScreenObject, winGoal int) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	menu := BaseObject{
		false,
		true,
		3,
		1,
		style,
		0,
		fmt.Sprintf("Meteorite Goal: %v", winGoal),
		make(chan (struct{})),
		make(chan (struct{})),
	}
	for {
		menuChan <- &menu
	}
}

func runMeteoriteCounter(menuChan chan ScreenObject) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	menu := BaseObject{
		false,
		true,
		3,
		2,
		style,
		0,
		fmt.Sprintf("Destroyed Meteorites: %v", destroyedMeteorites),
		make(chan (struct{})),
		make(chan (struct{})),
	}
	for {
		menu.View = fmt.Sprintf("Destroyed Meteorites: %v", destroyedMeteorites)
		menuChan <- &menu
	}
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
	newView := "Lifes: "
	for i := 0; i < counter.lifeNumber; i++ {
		newView += "♥"
		if i < counter.lifeNumber-1 {
			newView += " "
		}
	}
	counter.View = newView
}
