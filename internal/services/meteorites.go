package services

import "github.com/gdamore/tcell/v2"

func GenerateMeteorites(events chan ScreenObject) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		baseObject := BaseObject{
			events,
			false,
			true,
			float64(i),
			0,
			meteoriteStyle,
		}
		meteorite := Meteorite{baseObject, 0.01}
		go meteorite.Move()
	}
}

type Meteorite struct {
	BaseObject
	speed float64 // cells per iteration
}

func (this *Meteorite) Move() {
	for {
		if this.Active != true {
			break
		}
		if this.IsBlocked {
			continue
		}
		this.Y += this.speed
		this.IsBlocked = true
		this.Objects <- this
	}
}
