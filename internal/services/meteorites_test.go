package services_test

import (
	"context"
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
		true,
		0,
		initialY,
		tcell.StyleDefault.Background(tcell.ColorReset),
		speed,
		services.MeteoriteView1,
		make(chan struct{}),
		make(chan struct{}),
		true,
	}
	MockedScreenSvc := mocks.NewScreenSvc(t)
	MockedScreenSvc.On("GetScreenSize").Return(100, 100)
	meteorite := &services.Meteorite{
		BaseObject: baseObject,
		Objects:    chan<- services.ScreenObject(objectChannel),
		ScreenSvc:  MockedScreenSvc,
	}
	go meteorite.Move(context.Background())

	select {
	case receivedMeteorite := <-objectChannel:
		require.Equal(t, meteorite, receivedMeteorite)
	case <-time.After(1 * time.Second):
		t.Errorf("no meteorites in the object channel")
	}
	MockedScreenSvc.AssertNumberOfCalls(t, "GetScreenSize", 1)
	require.True(t, meteorite.Active)
	expectedY := initialY + meteorite.MaxSpeed
	require.Equal(t, expectedY, meteorite.Y)

	select {
	case <-objectChannel:
		t.Errorf("a blocked meteorite appears in the channel")
	case <-time.After(100 * time.Millisecond):
	}

	meteorite.Unblock()
	time.Sleep(100 * time.Millisecond)

	MockedScreenSvc.AssertNumberOfCalls(t, "GetScreenSize", 2)
	require.True(t, meteorite.Active)
	expectedY += meteorite.MaxSpeed
	require.Equal(t, expectedY, meteorite.Y)

	meteorite.Deactivate()
	time.Sleep(100 * time.Millisecond)
	require.False(t, meteorite.Active)
}

func TestMeteorite_MoveAndLeaveScreen(t *testing.T) {
	objectChannel := make(chan services.ScreenObject)
	initialY, speed := float64(2), float64(1)
	baseObject := services.BaseObject{
		false,
		true,
		0,
		initialY,
		tcell.StyleDefault.Background(tcell.ColorReset),
		speed,
		services.MeteoriteView1,
		make(chan struct{}),
		make(chan struct{}),
		true,
	}
	MockedScreenSvc := mocks.NewScreenSvc(t)
	MockedScreenSvc.On("GetScreenSize").Return(1, 1)
	meteorite := &services.Meteorite{
		BaseObject: baseObject,
		Objects:    chan<- services.ScreenObject(objectChannel),
		ScreenSvc:  MockedScreenSvc,
	}
	go meteorite.Move(context.Background())

	select {
	case <-objectChannel:
	case <-time.After(100 * time.Millisecond):
	}

	meteorite.Unblock()
	time.Sleep(100 * time.Millisecond)
	require.False(t, meteorite.Active)
}

func TestMeteorite_Collide(t *testing.T) {
	baseObject := services.BaseObject{
		false,
		true,
		0,
		0,
		tcell.StyleDefault.Background(tcell.ColorReset),
		1,
		services.MeteoriteView1,
		make(chan struct{}),
		make(chan struct{}),
		true,
	}
	objectChannel := make(chan<- services.ScreenObject)
	screenMock := mocks.NewScreenSvc(t)
	tests := []struct {
		name              string
		collisionObjects  []services.ScreenObject
		expectActiveState bool
	}{
		{
			"no objects collide",
			[]services.ScreenObject{},
			true,
		},
		{
			"collide with meoteorites",
			[]services.ScreenObject{
				&services.Meteorite{
					BaseObject: baseObject,
					Objects:    objectChannel,
					ScreenSvc:  screenMock,
				},
				&services.Meteorite{
					BaseObject: baseObject,
					Objects:    objectChannel,
					ScreenSvc:  screenMock,
				},
			},
			true,
		},
		{
			"collide with the other object",
			[]services.ScreenObject{
				&services.Meteorite{
					BaseObject: baseObject,
					Objects:    objectChannel,
					ScreenSvc:  screenMock,
				},
				&baseObject,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meteorite := &services.Meteorite{
				BaseObject: baseObject,
				Objects:    objectChannel,
				ScreenSvc:  screenMock,
			}
			meteorite.Collide(context.Background(), tt.collisionObjects)
			require.Equal(t, tt.expectActiveState, meteorite.Active)
		})
	}
}
