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
		head := styleFuncChain{
			style.Bold,
			true,
			nil,
		}
		next := styleFuncChain{
			style.Bold,
			false,
			nil,
		}
		next2 := styleFuncChain{
			style.Dim,
			true,
			nil,
		}
		next3 := styleFuncChain{
			style.Dim,
			false,
			nil,
		}
		next3.next = &head
		next2.next = &next3
		next.next = &next2
		head.next = &next
		star := Star{baseObject, starsChan, &head}
		go star.Blink()
		usedCoords[starLine][starColumn] = true
	}
}

type styleFuncChain struct {
	self func(bool) tcell.Style
	argument bool
	next *styleFuncChain
}

type Star struct {
	BaseObject
	StarChan   chan<- ScreenObject
	styleChain *styleFuncChain
}

func (star *Star) Blink() {
	tickOffset := time.Duration(rand.Intn(400) + 150)
	ticker := time.NewTicker(tickOffset * time.Millisecond)
	for {
		star.StarChan <- star
		select {
		case <-ticker.C:
			star.styleChain = star.styleChain.next
			star.Style = star.styleChain.self(star.styleChain.argument)
		default:
		}
	}
}
