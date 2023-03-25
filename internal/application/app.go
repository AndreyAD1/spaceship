package application

import (
	"log"
	"os"
	"time"

	"github.com/AndreyAD1/spaceship/internal/services"
)

type Application struct {
	Logger       *log.Logger
	FrameTimeout time.Duration
}

func GetApplication(logger *log.Logger) Application {
	return Application{logger, 400 * time.Millisecond}
}

func (this Application) quit(screenSvc services.ScreenService) {
	screenSvc.Finish()
	os.Exit(0)
}

func (this Application) Run() error {
	screenService, err := services.GetScreenService()
	if err != nil {
		return err
	}
	defer this.quit(screenService)

	objectChannel := make(chan *services.ScreenObject, 1)
	objectsToLoose := []*services.ScreenObject{}
	services.GenerateMeteorites(objectChannel)
	this.Logger.Println("start an event loop")
	for {
		this.Logger.Println("get screen event")
		// TODO fix a Ctrl+C exit
		userEvent := screenService.GetScreenEvent()
		if userEvent == services.Exit {
			this.Logger.Println("received an exit signal")
			break
		}
		screenService.ClearScreen()
	ObjectLoop:
		for {
			this.Logger.Printf("get object info. Object to loose: %v", objectsToLoose)
			select {
			case object := <-objectChannel:
				this.Logger.Printf("receive an object %v-%v", object.X, object.Y)
				screenService.Draw(object)
				objectsToLoose = append(objectsToLoose, object)
			default:
				this.Logger.Printf("channel is empty, object: %v", objectsToLoose)
				for _, object := range objectsToLoose {
					this.Logger.Printf("loose object %v-%v", object.X, object.Y)
					object.Block <- struct{}{}
				}
				objectsToLoose = objectsToLoose[:0]
				break ObjectLoop
			}
		}
		screenService.ShowScreen()
		// TODO think about different object speeds
		time.Sleep(this.FrameTimeout)
	}
	return nil
}
