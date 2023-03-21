package main

import (
	"fmt"
	"time"

	"github.com/AndreyAD1/spaceship/frame"
	"github.com/gdamore/tcell/v2"
)

func drawMeteorite(screen tcell.Screen) {
	width, height := screen.Size()
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	meteorColumn := width / 2
	for line := 0; line < height - 3; line++ {
		screen.SetContent(meteorColumn, line, 'O', nil, meteoriteStyle)
		screen.Show()
		time.Sleep(250 * time.Millisecond)
		screen.Clear()
	}
	screen.Show()
}

func main() {
	fmt.Println("a spaceship game will be here")
	newFrame := frame.DrawFrame()
	defer newFrame.Clear()
	drawMeteorite(newFrame)
}