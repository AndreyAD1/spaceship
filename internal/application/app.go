package application

import (
	"os"
	"time"

	"github.com/AndreyAD1/spaceship/internal/config"
	"github.com/AndreyAD1/spaceship/internal/services"
)

type Application struct {
	DebugMode bool
	FrameTimeout time.Duration
}

func GetApplication(config config.StartupConfig) Application {
	return Application{config.Debug, 400 * time.Millisecond}
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

	objectChannel := make(chan *services.ScreenObject)
	objectsToLoose := []*services.ScreenObject{}
	services.GenerateMeteorites(objectChannel)
	for {
		// TODO fix a Ctrl+C exit
		userEvent := screenService.GetScreenEvent()
		if userEvent == services.Exit {
			break
		}
		screenService.ClearScreen()
	ObjectLoop:
		for {
			select {
			case object := <-objectChannel:
				screenService.Draw(object)
				objectsToLoose = append(objectsToLoose, object)
			default:
				for _, object := range objectsToLoose {
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
