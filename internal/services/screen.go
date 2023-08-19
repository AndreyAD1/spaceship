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
	GoUp
	GoDown
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
	NewObjectList() [][][]ScreenObject
	GetScreenSize() (int, int)
}

type ScreenService struct {
	screen         tcell.Screen
	exitChannel    chan struct{}
	controlChannel chan ScreenEvent
	width          int
	height         int
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
	// Sometimes a screen appears only after a resizing
	screen.SetSize(width+1, height)
	newSvc := ScreenService{
		screen,
		make(chan struct{}),
		// a channel buffer allows a user to exit in a gameover state
		make(chan ScreenEvent, 50),
		width + 1,
		height,
	}
	return &newSvc, nil
}

func (screenSvc *ScreenService) PollScreenEvents(ctx context.Context) {
	logger := log.FromContext(ctx)
	events := make(chan tcell.Event)
	quit := make(chan struct{})
	go screenSvc.screen.ChannelEvents(events, quit)
	var event tcell.Event
	for {
		select {
		case <-ctx.Done():
			close(quit)
			return
		case event = <-events:
		}
		if ev, ok := event.(*tcell.EventKey); ok {
			logger.Debugf("process a key event %v", ev.Key())
			switch ev.Key() {
			case tcell.KeyLeft:
				screenSvc.controlChannel <- GoLeft
			case tcell.KeyRight:
				screenSvc.controlChannel <- GoRight
			case tcell.KeyUp:
				screenSvc.controlChannel <- GoUp
			case tcell.KeyDown:
				screenSvc.controlChannel <- GoDown
			case tcell.KeyRune:
				logger.Debugf("key \"%c\" is pressed", ev.Rune())
				if ev.Rune() == ' ' {
					screenSvc.controlChannel <- Shoot
				}
			case tcell.KeyEscape:
				close(screenSvc.exitChannel)
				close(quit)
				return
			case tcell.KeyCtrlC:
				close(screenSvc.exitChannel)
				close(quit)
				return
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
	width, height := screenSvc.GetScreenSize()
	roundX, roundY := int(math.Round(x)), int(math.Round(y))
	xIsOutside := roundX >= width-1 || roundX < 0
	yIsOutside := roundY >= height || roundY < 0
	return !(xIsOutside || yIsOutside)
}

func (screenSvc *ScreenService) Draw(obj ScreenObject) {
	coords, characters := obj.GetViewCoordinates()
	for i, character := range characters {
		x, y := coords[i][0], coords[i][1]
		if character == 0x85 || !unicode.IsSpace(character) {
			screenSvc.screen.SetContent(x, y, character, nil, obj.GetStyle())
		}
	}
}

func (screenSvc *ScreenService) NewObjectList() [][][]ScreenObject {
	width, height := screenSvc.GetScreenSize()
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
	return screenSvc.width, screenSvc.height
}
