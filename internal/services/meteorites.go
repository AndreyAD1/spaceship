package services

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

func GenerateMeteorites(events chan ScreenObject, screenSvc *ScreenService) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	width, _ := screenSvc.GetScreenSize()
	for {
		if rand.Float32() < 0.99 {
			continue
		}
		baseObject := BaseObject{
			false,
			false,
			true,
			float64(rand.Intn(width)),
			0,
			meteoriteStyle,
			0.01,
			"O",
		}
		meteorite := Meteorite{baseObject, events, screenSvc}
		go meteorite.Move()
		time.Sleep(time.Second)
	}
}

type Meteorite struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
}

func (this *Meteorite) Move() {
	for {
		if this.Active != true {
			break
		}
		if this.IsBlocked {
			continue
		}
		newY := this.Y + this.Speed
		if !this.ScreenSvc.IsInsideScreen(this.X, newY) {
			this.Deactivate()
			break
		}
		this.Y = newY
		this.IsBlocked = true
		this.Objects <- this
	}
}
