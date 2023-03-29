package application

import (
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

func (this Application) Run() error {
	screenService, err := services.GetScreenService()
	if err != nil {
		return err
	}
	defer screenService.Exit()

	interactiveObjects := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.GameOver)
	services.GenerateMeteorites(interactiveObjects, screenService)
	services.GenerateShip(interactiveObjects, screenService, gameoverChannel)
	go screenService.PollScreenEvents()

	screenObjects := screenService.GetObjectList() 
	this.Logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			break
		}
	ObjectLoop:
		for {
			select {
			case object := <-interactiveObjects:
				x, y := object.GetCoordinates()
				screenObjects[y][x] = append(screenObjects[y][x], object)
			default:
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
		select {
		case gameover := <- gameoverChannel:
			screenService.Draw(gameover)
		default:
		} 

		screenService.ShowScreen()
		time.Sleep(this.FrameTimeout)
		screenService.ClearScreen()
	}
	this.Logger.Debug("finish the event loop")
	return nil
}
