package services

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

const MeteoriteView = `  ___
 /   \
/     /
\____/`

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
			-7,
			meteoriteStyle,
			0.01,
			MeteoriteView,
		}
		meteorite := Meteorite{baseObject, events, screenSvc}
		go meteorite.Move()
		time.Sleep(time.Second * 1)
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
		_, height := this.ScreenSvc.GetScreenSize()
		if newY > float64(height) + 2 {
			this.Deactivate()
			break
		}
		this.Y = newY
		this.IsBlocked = true
		this.Objects <- this
	}
}
