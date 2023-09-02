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
) ([]ScreenObject, chan int, chan int) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
	template := "%v\nDestroyed Meteorites: %v/%v\nLifes: %v"
	initialView := fmt.Sprintf(
		template,
		levelName,
		destroyedMeteorites,
		winGoal,
		getLifeView(initialLifeNumber),
	)
	screenObject := BaseObject{
		false,
		true,
		3,
		1,
		style,
		0,
		initialView,
		make(chan (struct{})),
		make(chan (struct{})),
		false,
	}
	menu := Menu{
		screenObject,
		template,
		levelName,
		initialLifeNumber,
		winGoal,
		0,
	}
	lifeChannel := make(chan int)
	destroyedMeteoriteChannel := make(chan int)
	go menu.runMenu(ctx, lifeChannel, destroyedMeteoriteChannel)
	return []ScreenObject{&menu}, lifeChannel, destroyedMeteoriteChannel
}

type Menu struct {
	BaseObject
	viewTemplate        string
	levelName           string
	lifeNumber          int
	winGoal             int
	destroyedMeteorites int
}

func (menu *Menu) runMenu(
	ctx context.Context,
	lifeChannel,
	destroyedMeteoriteChannel <-chan int,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case lifeNumber := <-lifeChannel:
			menu.lifeNumber = lifeNumber
			menu.View = menu.getNewView()
		case meteoriteNumber := <-destroyedMeteoriteChannel:
			menu.destroyedMeteorites = meteoriteNumber
			menu.View = menu.getNewView()
		}
	}
}

func (menu *Menu) getNewView() string {
	newView := fmt.Sprintf(
		menu.viewTemplate,
		menu.levelName,
		menu.destroyedMeteorites,
		menu.winGoal,
		getLifeView(menu.lifeNumber),
	)
	return newView
}

func getLifeView(lifeNumber int) string {
	newView := ""
	for i := 0; i < lifeNumber; i++ {
		newView += "♥"
		if i < lifeNumber-1 {
			newView += " "
		}
	}
	return newView
}
