package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/AndreyAD1/spaceship/frame"
	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/caarlos0/env/v7"
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
		if screen.HasPendingEvent() {
			event := screen.PollEvent()
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					screen.Fini()
					os.Exit(0)
				}
			}
		}
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

func quit(screen tcell.Screen) {
	screen.Fini()
	os.Exit(0)
}

func main() {
	debug := flag.String("debug", "", "run in a debug mode")
	flag.Parse()
	configuration := config.StartupConfig{}
	err := env.Parse(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	if *debug == "true" || *debug == "false" {
		configuration.Debug = *debug == "true"
	}
	screen := frame.GetScreen()
	defer quit(screen)
	draw(screen)
}
