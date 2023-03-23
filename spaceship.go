package main

import (
	"time"

	"github.com/AndreyAD1/spaceship/frame"
	"github.com/gdamore/tcell/v2"
)

type ScreenObject struct {
	Events    chan<- *ScreenObject
	Block     chan struct{}
	Eliminate chan struct{}
	X         int
	Y         int
	Style     tcell.Style
}

func (this *ScreenObject) Move() {
	for {
		select {
		case <-this.Eliminate:
			break
		default:
		}
		this.Y++
		this.Events <- this
		<-this.Block
	}
}

func (this *ScreenObject) Draw(screen tcell.Screen) {
	width, height := screen.Size()
	if this.X > width || this.Y > height {
		this.Eliminate <- struct{}{}
		return
	}
	screen.SetContent(this.X, this.Y, 'O', nil, this.Style)
}

func generateMeteorites(events chan *ScreenObject) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	for i := 0; i < 10; i += 3 {
		meteorite := ScreenObject{
			events,
			make(chan struct{}),
			make(chan struct{}),
			i,
			0,
			meteoriteStyle,
		}
		go meteorite.Move()
	}
}

func draw(screen tcell.Screen) {
	objectChannel := make(chan *ScreenObject)
	objectsToLoose := []*ScreenObject{}
	generateMeteorites(objectChannel)
	for {
		screen.Clear()
	ObjectLoop:
		for {
			select {
			case object := <-objectChannel:
				object.Draw(screen)
				objectsToLoose = append(objectsToLoose, object)
			default:
				for _, object := range objectsToLoose {
					object.Block <- struct{}{}
				}
				objectsToLoose = objectsToLoose[:0]
				break ObjectLoop
			}
		}
		screen.Show()
		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	screen := frame.GetScreen()
	defer screen.Clear()
	draw(screen)
}
