package services

import "github.com/gdamore/tcell/v2"

type ScreenObject interface {
	GetCoordinates() (int, int)
	GetStyle() tcell.Style
	Move()
	Unblock()
	Deactivate()
}
