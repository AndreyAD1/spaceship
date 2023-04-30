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
	MockedScreenSvc.AssertNumberOfCalls(t, "GetScreenSize", 1)
	require.True(t, meteorite.IsBlocked)
	require.True(t, meteorite.Active)
	expectedY := initialY + meteorite.Speed
	require.Equal(t, expectedY, meteorite.Y)

	select {
	case <-objectChannel:
		t.Errorf("a blocked meteorite appears in the channel")
	case <-time.After(100 * time.Millisecond):
	}

	meteorite.IsBlocked = false

	select {
	case receivedMeteorite := <-objectChannel:
		require.Equal(t, meteorite, receivedMeteorite)
	case <-time.After(1 * time.Second):
		t.Errorf("no meteorites in the object channel")
	}
	MockedScreenSvc.AssertNumberOfCalls(t, "GetScreenSize", 2)
	require.True(t, meteorite.IsBlocked)
	require.True(t, meteorite.Active)
	expectedY += meteorite.Speed
	require.Equal(t, expectedY, meteorite.Y)

	meteorite.Active = false
	select {
	case <-objectChannel:
		t.Errorf("a deactivated meteorite appears in the channel")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestMeteorite_MoveAndLeaveScreen(t *testing.T) {
	objectChannel := make(chan services.ScreenObject)
	initialY, speed := float64(2), float64(1)
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
	MockedScreenSvc.On("GetScreenSize").Return(1, 1)
	meteorite := &services.Meteorite{
		BaseObject: baseObject,
		Objects:    chan<- services.ScreenObject(objectChannel),
		ScreenSvc:  MockedScreenSvc,
	}
	go meteorite.Move()
	select {
	case <- objectChannel:
	case <- time.After(100 * time.Millisecond):
	}
	meteorite.IsBlocked = false
	select {
	case receivedMeteorite := <-objectChannel:
		t.Errorf("a meteorite left the screen but appeared in the channel")
		require.Equal(t, meteorite, receivedMeteorite)
	case <-time.After(100 * time.Millisecond):
	}
	require.False(t, meteorite.Active)
}