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
	ScreenSvc ScreenSvc
}

func (meteorite *Meteorite) Move() {
	for {
		if meteorite.Active != true {
			break
		}
		if meteorite.IsBlocked {
			continue
		}
		newY := meteorite.Y + meteorite.Speed
		_, height := meteorite.ScreenSvc.GetScreenSize()
		if newY > float64(height)+2 {
			meteorite.Deactivate()
			break
		}
		meteorite.Y = newY
		meteorite.IsBlocked = true
		meteorite.Objects <- meteorite
	}
}

func (meteorite *Meteorite) Collide(objects []ScreenObject) {
	allObjectsAreMeteors := true
	for _, obj := range objects {
		switch obj.(type) {
		case *Meteorite:
		default:
			allObjectsAreMeteors = false
		}
	}
	if !allObjectsAreMeteors {
		meteorite.Deactivate()
	}
}
