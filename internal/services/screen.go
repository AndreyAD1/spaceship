package services

import "github.com/gdamore/tcell/v2"

type ScreenEvent int

const (
	NoEvent ScreenEvent = iota
	Exit
	GoLeft
	GoRight
)

type ScreenService struct {
	screen tcell.Screen
}

func GetScreenService() (ScreenService, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return ScreenService{}, err
	}
	if err := screen.Init(); err != nil {
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
			if ev.Key() == tcell.KeyLeft {
				return GoLeft
			}
			if ev.Key() == tcell.KeyRight {
				return GoRight
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

func (this ScreenService) Draw(obj ScreenObject) {
	width, height := this.screen.Size()
	x, y := obj.GetCoordinates()
	if x > width || y > height {
		obj.Deactivate()
		return
	}
	this.screen.SetContent(x, y, 'O', nil, obj.GetStyle())
}
