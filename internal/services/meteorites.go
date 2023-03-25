package services

import "github.com/gdamore/tcell/v2"

func GenerateMeteorites(events chan *ScreenObject) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		meteorite := ScreenObject{
			events,
			make(chan struct{}),
			make(chan struct{}),
			i,
			0,
			meteoriteStyle,
		}
		go meteorite.Move()
	}
}
