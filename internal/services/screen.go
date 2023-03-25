package services

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type ScreenEvent int

const (
	NoEvent ScreenEvent = iota
	Exit
)

type ScreenService struct {
	screen tcell.Screen
}

func GetScreenService() (ScreenService, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		// TODO add proper logs
		log.Printf("%+v", err)
		return ScreenService{}, err
	}
	if err := screen.Init(); err != nil {
		log.Printf("%+v", err)
		return ScreenService{}, err
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	defStyle = defStyle.Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)
	width, height := screen.Size()
	// Sometimes screen appears only after a resizing
	screen.SetSize(width+1, height)
	return ScreenService{screen}, nil
}

func (this ScreenService) GetScreenEvent() ScreenEvent {
	if this.screen.HasPendingEvent() {
		event := this.screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return Exit
			}
		}
	}
	return NoEvent
}

func (this ScreenService) ClearScreen() {
	this.screen.Clear()
}

func (this ScreenService) ShowScreen() {
	this.screen.Show()
}

func (this ScreenService) Finish() {
	this.screen.Fini()
}

func (this ScreenService) Draw(obj *ScreenObject) {
	width, height := this.screen.Size()
	if obj.X > width || obj.Y > height {
		obj.Eliminate <- struct{}{}
		return
	}
	this.screen.SetContent(obj.X, obj.Y, 'O', nil, obj.Style)
}
