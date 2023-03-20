package frame

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func DrawFrame() tcell.Screen {
	frame, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := frame.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	defStyle = defStyle.Foreground(tcell.ColorReset)
	frame.SetStyle(defStyle)
	frame.SetContent(0, 0, 'H', nil, defStyle)
	frame.SetContent(1, 0, 'i', nil, defStyle)
	frame.SetContent(2, 0, '!', nil, defStyle)
	return frame
}