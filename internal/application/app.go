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
	gameoverChannel := make(chan *services.BaseObject)
	services.GenerateMeteorites(interactiveObjects, screenService)
	services.GenerateShip(interactiveObjects, screenService, gameoverChannel)
	go screenService.PollScreenEvents()

	screenObjects := screenService.GetObjectList()
	interObjects := []services.ScreenObject{}
	this.Logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			break
		}
	ObjectLoop:
		for {
			select {
			case object := <-interactiveObjects:
				interObjects = append(interObjects, object)
				coordinates := object.GetCoordinates()
				this.Logger.Debugf("object %v coordinates: %v", object.GetView(), coordinates)
				for _, coord_pair := range coordinates {
					x, y := coord_pair[0], coord_pair[1]
					if screenService.IsInsideScreen(float64(x), float64(y)) {
						screenObjects[y][x] = append(screenObjects[y][x], object)
					}
				}
			default:
				break ObjectLoop
			}
		}
		this.Logger.Debugf("screen objects: %v", screenObjects)
		for y, row := range screenObjects {
			for x, objects := range row {
				if len(objects) == 0 {
					continue
				}
				if len(objects) == 1 && !objects[0].GetDrawStatus() {
					screenService.Draw(objects[0])
					objects[0].MarkDrawn()
					screenObjects[y][x] = []services.ScreenObject{}
					continue
				}
				// collision occurred
				if len(objects) > 1 {
					for _, object := range objects {
						object.Deactivate()
					}
				}
				screenObjects[y][x] = []services.ScreenObject{}
			}
		}
		for _, object := range interObjects {
			object.Unblock()
		}
		interObjects = []services.ScreenObject{}
		select {
		case gameover := <-gameoverChannel:
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
