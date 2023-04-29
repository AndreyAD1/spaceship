package services_test

import (
	"testing"
	"time"

	"github.com/AndreyAD1/spaceship/internal/mocks"
	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestMeteorite_MoveAndDeactivate(t *testing.T) {
	objectChannel := make(chan services.ScreenObject)
	initialY, speed := float64(0), float64(1)
	baseObject := services.BaseObject{
		false,
		false,
		true,
		0,
		initialY,
		tcell.StyleDefault.Background(tcell.ColorReset),
		speed,
		services.MeteoriteView,
	}
	MockedScreenSvc := mocks.NewScreenSvc(t)
	MockedScreenSvc.On("GetScreenSize").Return(100, 100)
	meteorite := &services.Meteorite{
		BaseObject: baseObject,
		Objects:    chan<- services.ScreenObject(objectChannel),
		ScreenSvc:  MockedScreenSvc,
	}
	go meteorite.Move()

	select {
	case receivedMeteorite := <-objectChannel:
		require.Equal(t, meteorite, receivedMeteorite)
	case <-time.After(1 * time.Second):
		t.Errorf("no meteorites in the object channel")
	}
	require.True(t, meteorite.IsBlocked)
	MockedScreenSvc.AssertCalled(t, "GetScreenSize")

	select {
	case receivedMeteorite := <-objectChannel:
		t.Errorf("a blocked meteorite appears in the channel")
		require.Equal(t, meteorite, receivedMeteorite)
	case <-time.After(100 * time.Millisecond):
	}

	require.True(t, meteorite.IsBlocked)
	require.True(t, meteorite.Active)
	expectedY := initialY + meteorite.Speed
	require.Equal(t, expectedY, meteorite.Y)
	MockedScreenSvc.AssertNumberOfCalls(t, "GetScreenSize", 1)
}
