package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
	"github.com/charmbracelet/log"
)

func TestScreenService_PollScreenEvents_Exit(t *testing.T) {
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
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	ctx := log.WithContext(context.Background(), logger)
	go screenSvc.PollScreenEvents(ctx)
	screenMock.InjectKey(tcell.KeyCtrlC, 'h', tcell.ModNone)
	select {
	case <- exitChannel:
	case <- time.After(100 * time.Millisecond):
		t.Errorf("Ctrl+C caused no exit")
	}
	select {
	case _, ok := <- exitChannel:
		require.Falsef(t, ok, "exit channel is not close")
	case <- time.After(100 * time.Millisecond):
		t.Errorf("exit channel is not close")
	}
}
