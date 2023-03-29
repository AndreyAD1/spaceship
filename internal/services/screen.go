package services

import (
	"math"
	"github.com/gdamore/tcell/v2"
)

type ScreenEvent int

const (
	NoEvent ScreenEvent = iota
	Exit
	GoLeft
	GoRight
)

type ScreenService struct {
	screen         tcell.Screen
	exitChannel    chan struct{}
	controlChannel chan ScreenEvent
}

func GetScreenService() (*ScreenService, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	defStyle = defStyle.Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)
	width, height := screen.Size()
	// Sometimes screen appears only after a resizing
	screen.SetSize(width+1, height)
	newSvc := ScreenService{
		screen,
		make(chan struct{}),
		make(chan ScreenEvent),
	}
	return &newSvc, nil
}

func (this *ScreenService) PollScreenEvents() {
MainLoop:
	for {
		var event tcell.Event
		for this.screen.HasPendingEvent() {
			event = this.screen.PollEvent()
			if ev, ok := event.(*tcell.EventKey); ok {
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					this.exitChannel <- struct{}{}
					close(this.exitChannel)
					break MainLoop
				}
			}
		}
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				this.exitChannel <- struct{}{}
				close(this.exitChannel)
				break MainLoop
			}
			if ev.Key() == tcell.KeyLeft {
				this.controlChannel <- GoLeft
			}
			if ev.Key() == tcell.KeyRight {
				this.controlChannel <- GoRight
			}
		}
	}
}

func (this *ScreenService) Exit() bool {
	select {
	case <-this.exitChannel:
		return true
	default:
		return false
	}
}

func (this *ScreenService) GetControlEvent() ScreenEvent {
	select {
	case event := <-this.controlChannel:
		return event
	default:
		return NoEvent
	}
}

func (this *ScreenService) ClearScreen() {
	this.screen.Clear()
}

func (this *ScreenService) ShowScreen() {
	this.screen.Show()
}

func (this *ScreenService) Finish() {
	this.screen.Fini()
}

func (this *ScreenService) IsInsideScreen(x, y float64) bool {
	width, height := this.screen.Size()
	roundX, roundY := int(math.Round(x)), int(math.Round(y))
	if roundX >= width - 1 || roundX < 0 || roundY >= height || roundY < 0 {
		return false
	}
	return true
}

func (this *ScreenService) Draw(obj ScreenObject) {
	x, y := obj.GetCoordinates()
	this.screen.SetContent(x, y, 'O', nil, obj.GetStyle())
}

func (this *ScreenService) GetObjectList() [][][]ScreenObject {
	width, height := this.screen.Size()
	newList := make([][][]ScreenObject, height)
	for i := 0; i < height; i++ {
		newList[i] = make([][]ScreenObject, width)
		for j := 0; j < width; j++ {
			newList[i][j] = []ScreenObject{}
		}
	}
	return newList
}
