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
		time.Sleep(time.Millisecond * 1500)
		if rand.Float32() < 0.4 {
			continue
		}
		baseObject := BaseObject{
			false,
			false,
			true,
			float64(rand.Intn(width - 2)),
			-6,
			meteoriteStyle,
			0.02,
			MeteoriteView,
		}
		meteorite := Meteorite{baseObject, events, screenSvc}
		go meteorite.Move()
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

func (this *Meteorite) Collide(objects []ScreenObject) {
	allObjectsAreMeteors := true
	for _, obj := range objects {
		switch obj.(type) {
		case *Meteorite:
		default:
			allObjectsAreMeteors = false
		}
	}
	if !allObjectsAreMeteors {
		this.Deactivate()
	}
}
