package services

import "github.com/gdamore/tcell/v2"

func GenerateMeteorites(events chan ScreenObject, sreencSvc *ScreenService) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		baseObject := BaseObject{
			events,
			sreencSvc,
			false,
			true,
			float64(i),
			0,
			meteoriteStyle,
			0.01,
		}
		meteorite := Meteorite{baseObject}
		go meteorite.Move()
	}
}

type Meteorite struct {
	BaseObject
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
