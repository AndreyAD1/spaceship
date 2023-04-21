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

func (app Application) Run() error {
	screenService, err := services.NewScreenService()
	if err != nil {
		return err
	}
	defer screenService.Exit()

	interactiveObjects := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.BaseObject)
	go services.GenerateMeteorites(interactiveObjects, screenService)
	services.GenerateShip(interactiveObjects, screenService, gameoverChannel)
	go screenService.PollScreenEvents()

	screenObjects := screenService.GetObjectList()
	interObjects := []services.ScreenObject{}
	app.Logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			break
		}
	ObjectLoop:
		for {
			select {
			case obj := <-interactiveObjects:
				interObjects = append(interObjects, obj)
				coordinates := obj.GetViewCoordinates()
				for _, coord_pair := range coordinates {
					x, y := coord_pair[0], coord_pair[1]
					if screenService.IsInsideScreen(float64(x), float64(y)) {
						screenObjects[y][x] = append(screenObjects[y][x], obj)
					}
				}
			default:
				break ObjectLoop
			}
		}
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
						object.Collide(objects)
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
		time.Sleep(app.FrameTimeout)
		screenService.ClearScreen()
	}
	app.Logger.Debug("finish the event loop")
	return nil
}
