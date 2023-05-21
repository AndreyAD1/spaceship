package application

import (
	"context"
	"time"

	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/charmbracelet/log"
)

type Application struct {
	Logger       *log.Logger
	FrameTimeout time.Duration
}

func NewApplication(logger *log.Logger) Application {
	return Application{logger, 10 * time.Millisecond}
}

func (app Application) Run() error {
	ctx := log.WithContext(context.Background(), app.Logger)
	screenService, err := services.NewScreenService()
	if err != nil {
		return err
	}
	defer screenService.Finish()

	// starChannel := make(chan services.ScreenObject)
	interactiveChannel := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.BaseObject)

	// services.GenerateStars(starChannel, screenService)
	go services.GenerateMeteorites(interactiveChannel, screenService)
	services.GenerateShip(interactiveChannel, screenService, gameoverChannel)
	go screenService.PollScreenEvents(ctx)

	app.Logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			break
		}
		// drawStars(starChannel, screenService)
		processInteractiveObjects(interactiveChannel, screenService)
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

func drawStars(starChan chan services.ScreenObject, screenSvc *services.ScreenService) {
	for {
		select {
		case star := <- starChan:
			screenSvc.Draw(star)
		default:
			return
		}
	}
}

func processInteractiveObjects(
	objectChannel chan services.ScreenObject,
	screenService *services.ScreenService,
) {
	screenObjects, interObjects := getScreenObjects(objectChannel, screenService)
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
		}
	}

	for _, object := range interObjects {
		object.Unblock()
	}
}

func getScreenObjects(
	objectChannel chan services.ScreenObject,
	screenService *services.ScreenService,
) ([][][]services.ScreenObject, []services.ScreenObject) {
	screenObjects := screenService.NewObjectList()
	interObjects := []services.ScreenObject{}
	for {
		select {
		case obj := <-objectChannel:
			interObjects = append(interObjects, obj)
			coordinates, _ := obj.GetViewCoordinates()
			for _, coord_pair := range coordinates {
				x, y := coord_pair[0], coord_pair[1]
				if screenService.IsInsideScreen(float64(x), float64(y)) {
					screenObjects[y][x] = append(screenObjects[y][x], obj)
				}
			}
		default:
			return screenObjects, interObjects
		}
	}
}
