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
		if shell.Active != true {
			break
		}
		if shell.IsBlocked {
			continue
		}
		newY := shell.Y - shell.Speed
		if !shell.ScreenSvc.IsInsideScreen(shell.X, newY) {
			shell.Deactivate()
			break
		}
		shell.Y = newY
		shell.IsBlocked = true
		shell.Objects <- shell
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
