package services

import "github.com/gdamore/tcell/v2"

func GenerateMeteorites(events chan *ScreenObject) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		meteorite := ScreenObject{
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
