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
	return Application{logger, 10 * time.Millisecond}
}

func (this Application) quit(screenSvc *services.ScreenService) {
	screenSvc.Finish()
	os.Exit(0)
}

func (this Application) Run() error {
	screenService, err := services.GetScreenService()
	if err != nil {
		return err
	}
	defer this.quit(screenService)

	objectChannel := make(chan services.ScreenObject)
	services.GenerateMeteorites(objectChannel, screenService)
	services.GenerateShip(screenService, objectChannel)
	go screenService.PollScreenEvents()

	screenObjects := screenService.GetObjectList() 
	this.Logger.Debug("start an event loop")
	for {
		this.Logger.Debug("main event loop")
		if screenService.Exit() {
			break
		}
	ObjectLoop:
		for {
			// this.Logger.Debugf("get object info. Objects: %v", screenObjects)
			select {
			case object := <-objectChannel:
				x, y := object.GetCoordinates()
				if y < len(screenObjects) && x < len(screenObjects[y]) {
					screenObjects[y][x] = append(screenObjects[y][x], object)
				}
			default:
				// this.Logger.Debugf("channel is empty, objects: %v", screenObjects)
				break ObjectLoop
			}
		}
		for y, row := range screenObjects {
			for x, objects := range row {
				if len(objects) == 0 {
					continue
				}
				if len(objects) == 1 {
					screenService.Draw(objects[0])
					objects[0].Unblock()
					screenObjects[y][x] = []services.ScreenObject{}
					continue
				}
				// collision occurred
				for _, object := range objects {
					object.Deactivate()
					object.Unblock()
				}
				screenObjects[y][x] = []services.ScreenObject{}
			}
		}

		screenService.ShowScreen()
		time.Sleep(this.FrameTimeout)
		screenService.ClearScreen()
	}
	this.Logger.Debug("finish the event loop")
	return nil
}
