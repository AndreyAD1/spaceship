package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

func GenerateStars(
	ctx context.Context, 
	screenSvc *ScreenService,
) []ScreenObject {
	width, height := screenSvc.GetScreenSize()
	screenSquare := width * height
	usedCoords := make([][]bool, height)
	for i := range usedCoords {
		usedCoords[i] = make([]bool, width)
	}
	stars := []ScreenObject{}
	for i := 0; i < screenSquare/15; i++ {
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
			true,
			float64(starColumn),
			float64(starLine),
			style,
			0,
			"*",
			make(chan (struct{})),
			make(chan (struct{})),
			false,
		}
		star := Star{baseObject}
		go star.Blink(ctx)
		usedCoords[starLine][starColumn] = true
		stars = append(stars, &star)
	}
	return stars
}

type Star struct {
	BaseObject
}

func (star *Star) Blink(ctx context.Context) {
	tickOffset := time.Duration(rand.Intn(400) + 200)
	ticker := time.NewTicker(tickOffset * time.Millisecond)
	defer ticker.Stop()
	tickPhase := 0
	for {
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

		select {
		case <-ctx.Done():
			return
		case <-star.UnblockCh:
		}
	}
}
