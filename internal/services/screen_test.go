package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestScreenService_PollScreenEvents_Exit(t *testing.T) {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	ctx := log.WithContext(context.Background(), logger)
	commands := []tcell.Key{
		tcell.KeyLeft,
		tcell.KeyLeft,
		tcell.KeyRune,
		tcell.KeyRight,
	}

	tests := []struct {
		name       string
		pushedKeys []tcell.Key
	}{
		{"immediate exit Ctrl+C", []tcell.Key{tcell.KeyCtrlC}},
		{"immediate exit Escape", []tcell.Key{tcell.KeyEscape}},
		{
			"discard other keys and exit after Ctrl+C",
			append(commands, tcell.KeyCtrlC),
		},
		{
			"discard other keys and exit after Escape",
			append(commands, tcell.KeyEscape),
		},
		{
			"several exit commands",
			[]tcell.Key{tcell.KeyCtrlC, tcell.KeyEscape},
		},
		{
			"several Ctrl+C commands",
			[]tcell.Key{tcell.KeyCtrlC, tcell.KeyCtrlC},
		},
		{
			"an exit command is in the middle",
			append(commands, []tcell.Key{tcell.KeyEscape, tcell.KeyRune}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenMock := tcell.NewSimulationScreen("")
			defer screenMock.Fini()
			err := screenMock.Init()
			require.NoError(t, err)
			exitChannel := make(chan struct{})
			screenSvc := &ScreenService{
				screen:         screenMock,
				exitChannel:    exitChannel,
				controlChannel: make(chan ScreenEvent),
			}
			go screenSvc.PollScreenEvents(ctx)
			for _, key := range tt.pushedKeys {
				screenMock.InjectKey(key, ' ', tcell.ModNone)
			}
			select {
			case <-exitChannel:
			case <-time.After(10 * time.Millisecond):
				t.Errorf("Ctrl+C caused no exit")
			}
			select {
			case _, ok := <-exitChannel:
				require.Falsef(t, ok, "exit channel is not close")
			case <-time.After(10 * time.Millisecond):
				t.Errorf("exit channel is not close")
			}
		})
	}
}

func TestScreenService_PollScreenEvents_Controls(t *testing.T) {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	ctx := log.WithContext(context.Background(), logger)
	tests := []struct {
		name          string
		pushedKeys    []tcell.Key
		expectedEvent ScreenEvent
	}{
		{
			"one key event",
			[]tcell.Key{tcell.KeyRune},
			Shoot,
		},
		{
			"several events",
			[]tcell.Key{
				tcell.KeyLeft,
				tcell.KeyLeft,
				tcell.KeyRune,
				tcell.KeyRight,
			},
			GoRight,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenMock := tcell.NewSimulationScreen("")
			defer screenMock.Fini()
			err := screenMock.Init()
			require.NoError(t, err)

			exitChannel := make(chan struct{})
			screenSvc := &ScreenService{
				screen:         screenMock,
				exitChannel:    exitChannel,
				controlChannel: make(chan ScreenEvent),
			}
			go screenSvc.PollScreenEvents(ctx)

			for _, key := range tt.pushedKeys {
				screenMock.InjectKey(key, ' ', tcell.ModNone)
			}
			select {
			case event := <-screenSvc.controlChannel:
				require.Equal(t, tt.expectedEvent, event)
			case <-time.After(10 * time.Millisecond):
				t.Errorf("no control event")
			}

			screenMock.InjectKey(tcell.KeyRune, ' ', tcell.ModNone)
			select {
			case event := <-screenSvc.controlChannel:
				require.Equal(t, Shoot, event)
			case <-time.After(10 * time.Millisecond):
				t.Errorf("no control event")
			}

			select {
			case <-exitChannel:
				t.Errorf("channel is unexpectedly closed")
			case <-time.After(10 * time.Millisecond):
			}
		})
	}
}

func TestScreenService_IsInsideScreen(t *testing.T) {
	screenSize := 10
	tests := []struct {
		name     string
		x        float64
		y        float64
		expected bool
	}{
		{
			"inside a screen",
			float64(screenSize - 2),
			float64(screenSize - 1),
			true,
		},
		{
			"negative x",
			float64(-1),
			float64(screenSize - 2),
			false,
		},
		{
			"boundary x",
			float64(screenSize - 1),
			float64(screenSize - 2),
			false,
		},
		{
			"negative y",
			float64(screenSize - 2),
			float64(-1),
			false,
		},
		{
			"boundary y",
			float64(screenSize - 2),
			float64(screenSize),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenMock := tcell.NewSimulationScreen("")
			defer screenMock.Fini()
			err := screenMock.Init()
			require.NoError(t, err)
			screenSvc := &ScreenService{
				screen:         screenMock,
				exitChannel:    make(chan struct{}),
				controlChannel: make(chan ScreenEvent),
				width: screenSize,
				height: screenSize,
			}
			require.Equal(t, tt.expected, screenSvc.IsInsideScreen(tt.x, tt.y))
		})
	}
}
