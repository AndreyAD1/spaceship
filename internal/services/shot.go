package services

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

func Shot(
	ctx context.Context,
	screenSvc *ScreenService,
	objects chan<- ScreenObject,
	x, y float64,
) {
	baseObject := BaseObject{
		false,
		true,
		x,
		y,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.5,
		"|",
		make(chan struct{}),
		make(chan struct{}),
		true,
	}
	shot := Shell{baseObject, objects, screenSvc}
	go shot.Move(ctx)
}

type Shell struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
}

func (shell *Shell) Move(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case shell.Objects <- shell:
	}
	for {
		newY := shell.Y - shell.MaxSpeed
		if !shell.ScreenSvc.IsInsideScreen(shell.X, newY) {
			shell.Active = false
			break
		}
		shell.Y = newY
		select {
		case <-shell.Cancel:
			shell.Active = false
			return
		case <-shell.UnblockCh:
		}
	}
}

func (shell *Shell) Collide(ctx context.Context, objects []ScreenObject) bool {
	collisionWithAnotherShell := false
Loop:
	for _, obj := range objects {
		switch obj.(type) {
		case *Shell:
			if obj != shell && !obj.IsActive() {
				collisionWithAnotherShell = true
				break Loop
			}
		default:
			collisionWithAnotherShell = false
		}
	}
	if !collisionWithAnotherShell {
		shell.Deactivate()
		return true
	}
	return false
}
