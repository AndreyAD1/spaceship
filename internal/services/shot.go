package services

import "github.com/gdamore/tcell/v2"

func Shot(screenSvc *ScreenService, objects chan<- ScreenObject, x, y float64) {
	baseObject := BaseObject{
		false,
		false,
		true,
		x,
		y,
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.1,
		"|",
		make(chan struct{}),
		make(chan struct{}),
	}
	shot := Shell{baseObject, objects, screenSvc}
	go shot.Move()
}

type Shell struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc *ScreenService
}

func (shell *Shell) Move() {
	for {
		newY := shell.Y - shell.Speed
		if !shell.ScreenSvc.IsInsideScreen(shell.X, newY) {
			shell.Active = false
			break
		}
		shell.Y = newY
		shell.IsBlocked = true
		shell.Objects <- shell
		select {
		case <-shell.Cancel:
			return
		case <-shell.UnblockCh:
		}
	}
}

func (shell *Shell) Collide(objects []ScreenObject) {
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
	}
}
