package services

import (
	"math/rand"
	"time"

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
		style := tcell.StyleDefault.Background(tcell.ColorReset).Normal()
		baseObject := BaseObject{
			false,
			false,
			true,
			float64(starColumn),
			float64(starLine),
			style,
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
	tickOffset := time.Duration(rand.Intn(400) + 150)
	ticker := time.NewTicker(tickOffset * time.Millisecond)
	tickPhase := 0
	for {
		star.StarChan <- star
		select {
		case <-ticker.C:
			switch tickPhase {
			case 0: 
				star.Style = star.Style.Bold(true)
			case 1:
				star.Style = star.Style.Normal()
			case 2:
				star.Style = star.Style.Dim(true)
			case 3:
				star.Style = star.Style.Normal()
			}
			tickPhase++
			if tickPhase > 3 {
				tickPhase = 0
			}
		default:
		}
	}
}
