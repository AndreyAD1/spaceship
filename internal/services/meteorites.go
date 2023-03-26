package services

import "github.com/gdamore/tcell/v2"

func GenerateMeteorites(events chan ScreenObject) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		meteorite := Meteorite{
			events,
			false,
			true,
			i,
			0,
			meteoriteStyle,
		}
		go meteorite.Move()
	}
}

type Meteorite struct {
	Events    chan<- ScreenObject
	IsBlocked bool
	Active    bool
	X         int
	Y         int
	Style     tcell.Style
}

func (this *Meteorite) Move() {
	for {
		if this.Active != true {
			break
		}
		if this.IsBlocked {
			continue
		}
		this.Y++
		this.IsBlocked = true
		this.Events <- this
	}
}

func (this *Meteorite) Deactivate() {
	this.Active = false
}

func (this *Meteorite) Unblock() {
	this.IsBlocked = false
}

func (this *Meteorite) GetCoordinates() (int, int) {
	return this.X, this.Y
}

func (this *Meteorite) GetStyle() tcell.Style {
	return this.Style
}
