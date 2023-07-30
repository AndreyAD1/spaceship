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

const meteoriteGoal = 5
const shipLifes = 3

func NewApplication(logger *log.Logger) Application {
	return Application{logger, 10 * time.Millisecond}
}

func (app Application) Run(is_last_level bool) error {
	ctx := log.WithContext(context.Background(), app.Logger)
	screenService, err := services.NewScreenService()
	if err != nil {
		return err
	}
	defer screenService.Finish()

	menuChannel := make(chan services.ScreenObject)
	starChannel := make(chan services.ScreenObject)
	interactiveChannel := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.BaseObject)
	lifeChannel := services.GenerateMenu(menuChannel, meteoriteGoal)
	invulnerableChannel := make(chan services.ScreenObject)
	goalAchievedChannel := make(chan bool, 2)

	services.GenerateStars(starChannel, screenService)
	go services.GenerateMeteorites(
		interactiveChannel,
		invulnerableChannel,
		screenService,
		app.Logger,
	)
	services.GenerateShip(
		interactiveChannel,
		screenService,
		goalAchievedChannel,
		lifeChannel,
		invulnerableChannel,
		meteoriteGoal,
	)
	go screenService.PollScreenEvents(ctx)
	shipCollisions, meteoriteCollisions := 0, 0
	gameIsOver := false

	app.Logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			break
		}
		processInvulnerableObjects(starChannel, screenService)
		shipCollisions, meteoriteCollisions = processInteractiveObjects(
			interactiveChannel, 
			screenService, 
			shipCollisions, 
			meteoriteCollisions,
		)
		processInvulnerableObjects(invulnerableChannel, screenService)

		if shipCollisions >= shipLifes && !gameIsOver {
			go services.DrawGameOver(gameoverChannel, screenService)
			gameIsOver = true
		}
		if meteoriteCollisions >= meteoriteGoal && !gameIsOver {
			go services.DrawWin(gameoverChannel, screenService)
			gameIsOver = true
		}

		select {
		case gameover := <-gameoverChannel:
			screenService.Draw(gameover)
		default:
		}
		processInvulnerableObjects(menuChannel, screenService)
		screenService.ShowScreen()
		time.Sleep(app.FrameTimeout)
		screenService.ClearScreen()
	}
	app.Logger.Debug("finish the event loop")
	return nil
}

func processInvulnerableObjects(
	invulnerableChan chan services.ScreenObject,
	screenSvc *services.ScreenService,
) {
	invulnerableObjects := []services.ScreenObject{}
	for {
		select {
		case obj := <-invulnerableChan:
			screenSvc.Draw(obj)
			invulnerableObjects = append(invulnerableObjects, obj)
		default:
			for _, o := range invulnerableObjects {
				o.Unblock()
			}
			return
		}
	}
}

func processInteractiveObjects(
	objectChannel chan services.ScreenObject,
	screenService *services.ScreenService,
	spaceshipCollisions, destroyedMeteorites int,
) (int, int) {
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
					is_collided := object.Collide(objects)
					if is_collided {
						switch object.(type) {
						case *services.Spaceship:
							spaceshipCollisions++
						case *services.Meteorite:
							destroyedMeteorites++
						}
					}
				}
			}
		}
	}
	
	for _, object := range interObjects {
		if object.IsActive() {
			object.Unblock()
		}
	}
	return spaceshipCollisions, destroyedMeteorites
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
