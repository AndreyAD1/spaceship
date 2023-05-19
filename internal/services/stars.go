package services

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

func GenerateStars(starsChan chan ScreenObject, screenSvc *ScreenService) {
	width, height := screenSvc.GetScreenSize()
	screenSquare := width * height
	usedCoords := make([][]bool, height)
	for i := range usedCoords {
		usedCoords[i] = make([]bool, width)
	}
	for i := 0; i < screenSquare / 15; i++ {
		var starLine, starColumn int
		for {
			starLine, starColumn = rand.Intn(height), rand.Intn(width)
			if !usedCoords[starLine][starColumn] {
				break
			}
		}
		baseObject := BaseObject{
			false,
			false,
			true,
			float64(starColumn),
			float64(starLine),
			tcell.StyleDefault.Background(tcell.ColorReset),
			0,
			"*",
		}
		star := Star{baseObject, starsChan}
		go star.Blink()
		usedCoords[starLine][starColumn] = true
	}
}

type Star struct {
	BaseObject
	StarChan   chan<- ScreenObject
}

func (star *Star) Blink() {
	for {
		star.StarChan <- star
	}
}
