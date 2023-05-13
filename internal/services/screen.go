package services

import (
	"context"
	"fmt"
	"math"
	"unicode"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
)

type ScreenEvent int

const (
	NoEvent ScreenEvent = iota
	Exit
	GoLeft
	GoRight
	Shoot
)

type ScreenSvc interface {
	PollScreenEvents(ctx context.Context)
	Exit() bool
	GetControlEvent() ScreenEvent
	ClearScreen()
	ShowScreen()
	Finish()
	IsInsideScreen(x, y float64) bool
	Draw(obj ScreenObject)
	GetObjectList() [][][]ScreenObject
	GetScreenSize() (int, int)
}

type ScreenService struct {
	screen         tcell.Screen
	exitChannel    chan struct{}
	controlChannel chan ScreenEvent
}

func NewScreenService() (*ScreenService, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("can not get a new screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		err = fmt.Errorf(
			"can not initialize the new screen %v: %w",
			screen,
			err,
		)
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
		// a channel buffer allows a user to exit in a gameover state
		make(chan ScreenEvent, 15),
	}
	return &newSvc, nil
}

func (screenSvc *ScreenService) PollScreenEvents(ctx context.Context) {
	logger := log.FromContext(ctx)
MainLoop:
	for {
		var event tcell.Event
		for screenSvc.screen.HasPendingEvent() {
			event = screenSvc.screen.PollEvent()
			if ev, ok := event.(*tcell.EventKey); ok {
				logger.Debugf("receive a key event %v", ev.Key())
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					screenSvc.exitChannel <- struct{}{}
					close(screenSvc.exitChannel)
					break MainLoop
				}
			}
		}
		if ev, ok := event.(*tcell.EventKey); ok {
			logger.Debugf("process a key event %v", ev.Key())
			switch ev.Key() {
			case tcell.KeyLeft:
				logger.Debug("Left is pressed")
				screenSvc.controlChannel <- GoLeft
			case tcell.KeyRight:
				screenSvc.controlChannel <- GoRight
			case tcell.KeyRune:
				logger.Debugf("key \"%c\" is pressed", ev.Rune())
				if ev.Rune() == ' ' {
					screenSvc.controlChannel <- Shoot
				}
			}
		}
	}
}

func (screenSvc *ScreenService) Exit() bool {
	select {
	case <-screenSvc.exitChannel:
		return true
	default:
		return false
	}
}

func (screenSvc *ScreenService) GetControlEvent() ScreenEvent {
	select {
	case event := <-screenSvc.controlChannel:
		return event
	default:
		return NoEvent
	}
}

func (screenSvc *ScreenService) ClearScreen() {
	screenSvc.screen.Clear()
}

func (screenSvc *ScreenService) ShowScreen() {
	screenSvc.screen.Show()
}

func (screenSvc *ScreenService) Finish() {
	screenSvc.screen.Fini()
}

func (screenSvc *ScreenService) IsInsideScreen(x, y float64) bool {
	width, height := screenSvc.screen.Size()
	roundX, roundY := int(math.Round(x)), int(math.Round(y))
	if roundX >= width-1 || roundX < 0 || roundY >= height || roundY < 0 {
		return false
	}
	return true
}

func (screenSvc *ScreenService) Draw(obj ScreenObject) {
	initialX, y := obj.GetCornerCoordinates()
	view := obj.GetView()
	x := initialX
	for _, char := range view {
		if char == '\n' {
			y++
			x = initialX
			continue
		}
		if !unicode.IsSpace(char) {
			screenSvc.screen.SetContent(x, y, char, nil, obj.GetStyle())
		}
		x++
	}
}

func (screenSvc *ScreenService) GetObjectList() [][][]ScreenObject {
	width, height := screenSvc.screen.Size()
	newList := make([][][]ScreenObject, height)
	for i := 0; i < height; i++ {
		newList[i] = make([][]ScreenObject, width)
		for j := 0; j < width; j++ {
			newList[i][j] = []ScreenObject{}
		}
	}
	return newList
}

func (screenSvc *ScreenService) GetScreenSize() (int, int) {
	return screenSvc.screen.Size()
}
