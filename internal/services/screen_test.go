package services

import (
	"context"
	"encoding/binary"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestScreenService_PollScreenEvents_Exit(t *testing.T) {
	screenMock := tcell.NewSimulationScreen("")
	defer screenMock.Fini()
	err := screenMock.Init()
	require.NoError(t, err)

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	ctx := log.WithContext(context.Background(), logger)
	keySet := []byte{}
	keys := []uint64{
		uint64(tcell.KeyLeft),
		uint64(tcell.KeyLeft),
		uint64(' '),
		uint64(tcell.KeyRight),
	}
	for _, key := range keys {
		binary.AppendUvarint(keySet, key)
	}

	tests := []struct {
		name     string
		keyBytes []byte
	}{
		{"immediate exit", []byte{byte(tcell.KeyCtrlC)}},
		{"immediate exit", []byte{byte(tcell.KeyEscape)}},
		{
			"discard other keys and exit after Ctrc+C",
			append(keySet, byte(tcell.KeyCtrlC)),
		},
		{
			"discard other keys and exit after Escape",
			append(keySet, byte(tcell.KeyEscape)),
		},
		{
			"several exit commands",
			[]byte{byte(tcell.KeyCtrlC), byte(tcell.KeyEscape)},
		},
		{
			"several exit Ctrl+C",
			[]byte{byte(tcell.KeyCtrlC), byte(tcell.KeyCtrlC)},
		},
		{
			"exit command is in the middle",
			append(keySet, []byte{byte(tcell.KeyEscape), byte(' ')}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitChannel := make(chan struct{})
			screenSvc := &ScreenService{
				screen:         screenMock,
				exitChannel:    exitChannel,
				controlChannel: make(chan ScreenEvent),
			}
			go screenSvc.PollScreenEvents(ctx)
			screenMock.InjectKey(tcell.KeyCtrlC, 'h', tcell.ModNone)
			screenMock.InjectKeyBytes(tt.keyBytes)
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
