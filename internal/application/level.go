package application

import (
	"context"
	"fmt"
	"time"

	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/charmbracelet/log"
)

type levelConfig struct {
	name          string
	meteoriteGoal int
	lifes         int
	isLastLevel   bool
}

type level struct {
	name          string
	meteoriteGoal int
	lifes         int
	isLastLevel   bool
	frameTimeout  time.Duration
}

func NewLevel(config levelConfig, frameTimeout time.Duration) level {
	newLevel := level{
		config.name,
		config.meteoriteGoal,
		config.lifes,
		config.isLastLevel,
		frameTimeout,
	}
	return newLevel
}

func (lev level) Run(
	ctx context.Context,
	screenService *services.ScreenService,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	menuChannel := make(chan services.ScreenObject)
	starChannel := make(chan services.ScreenObject)
	interactiveChannel := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.BaseObject)
	invulnerableChannel := make(chan services.ScreenObject)
	var levelEnd <-chan time.Time

	lifeChannel := services.GenerateMenu(
		ctx,
		menuChannel,
		lev.name,
		lev.lifes,
		lev.meteoriteGoal,
	)
	services.GenerateStars(ctx, starChannel, screenService)
	go services.GenerateMeteorites(
		ctx,
		interactiveChannel,
		invulnerableChannel,
		screenService,
	)
	services.GenerateShip(
		ctx,
		interactiveChannel,
		screenService,
		lifeChannel,
		invulnerableChannel,
		lev.lifes,
	)
	go screenService.PollScreenEvents(ctx)
	shipCollisions, meteoriteCollisions := 0, 0
	gameIsOver := false

	logger := log.FromContext(ctx)
	logger.Debug("start an event loop")
	for {
		if screenService.Exit() {
			return fmt.Errorf("a user has stopped the game")
		}
		processInvulnerableObjects(starChannel, screenService)
		shipCollisions, meteoriteCollisions = processInteractiveObjects(
			ctx,
			interactiveChannel,
			screenService,
			shipCollisions,
			meteoriteCollisions,
		)
		processInvulnerableObjects(invulnerableChannel, screenService)

		if shipCollisions >= lev.lifes && !gameIsOver {
			go services.DrawLabel(ctx, gameoverChannel, screenService, services.GameOver)
			gameIsOver = true
		}
		if meteoriteCollisions >= lev.meteoriteGoal && !gameIsOver {
			if lev.isLastLevel {
				go services.DrawLabel(ctx, gameoverChannel, screenService, services.Win)
			}
			if !lev.isLastLevel && levelEnd == nil {
				go services.DrawLabel(ctx, gameoverChannel, screenService, services.Next)
				levelEnd = time.After(2 * time.Second)
			}
			gameIsOver = true
		}

		select {
		case gameover := <-gameoverChannel:
			screenService.Draw(gameover)
		case <-levelEnd:
			return nil
		default:
		}
		processInvulnerableObjects(menuChannel, screenService)
		screenService.ShowScreen()
		time.Sleep(lev.frameTimeout)
		screenService.ClearScreen()
	}
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
	ctx context.Context,
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
					is_collided := object.Collide(ctx, objects)
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
