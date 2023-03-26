package application

import (
	"os"
	"time"

	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/charmbracelet/log"
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
	this.Logger.Debug("start an event loop")
	for {
		this.Logger.Debug("get screen event")
		userEvent := screenService.GetScreenEvent()
		if userEvent == services.Exit {
			this.Logger.Debug("received an exit signal")
			break
		}
		screenService.ClearScreen()
	ObjectLoop:
		for {
			this.Logger.Debugf("get object info. Objects to loose: %v", objectsToLoose)
			select {
			case object := <-objectChannel:
				this.Logger.Debugf("receive an object %v-%v", object.X, object.Y)
				screenService.Draw(object)
				objectsToLoose = append(objectsToLoose, object)
			default:
				this.Logger.Debugf("channel is empty, object: %v", objectsToLoose)
				for _, object := range objectsToLoose {
					object.IsBlocked = false
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
