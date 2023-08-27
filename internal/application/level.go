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
var interactiveObjects map[services.ScreenObject]bool

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
	interactiveChannel := make(chan services.ScreenObject)
	gameoverChannel := make(chan *services.BaseObject)
	var levelEnd <-chan time.Time

	lifeChannel := services.GenerateMenu(
		ctx,
		menuChannel,
		lev.name,
		lev.lifes,
		lev.meteoriteGoal,
	)
	stars := services.GenerateStars(ctx, screenService)
	go services.GenerateMeteorites(
		ctx,
		interactiveChannel,
		screenService,
	)
	services.GenerateShip(
		ctx,
		interactiveChannel,
		screenService,
		lifeChannel,
		lev.lifes,
	)
	go screenService.PollScreenEvents(ctx)
	shipCollisions, meteoriteCollisions := 0, 0
	gameIsOver := false

	logger := log.FromContext(ctx)
	logger.Debug("start an event loop")
	interactiveObjects = make(map[services.ScreenObject]bool)
	for {
		if screenService.Exit() {
			return fmt.Errorf("a user has stopped the game")
		}
		processStaticObjects(stars, screenService)
		shipCollisions, meteoriteCollisions = processInteractiveObjects2(
			ctx,
			interactiveChannel,
			screenService,
			shipCollisions,
			meteoriteCollisions,
		)

		if shipCollisions >= lev.lifes && !gameIsOver {
			go services.DrawLabel(ctx, gameoverChannel, screenService, services.GameOver)
			gameIsOver = true
		}
		if meteoriteCollisions >= lev.meteoriteGoal && !gameIsOver {
			if lev.isLastLevel {
				go services.DrawLabel(ctx, gameoverChannel, screenService, services.Win)
			}
			if !lev.isLastLevel && levelEnd == nil {
				go services.DrawLabel(
					ctx,
					gameoverChannel,
					screenService,
					services.NextLevel,
				)
				levelEnd = time.After(2 * time.Second)
			}
			gameIsOver = true
		}

		select {
		case gameover := <-gameoverChannel:
			screenService.Draw(gameover)
		case <-levelEnd:
			logger.Debugf("finish the level")
			return nil
		default:
		}
		processInvulnerableObjects(menuChannel, screenService)
		screenService.ShowScreen()
		time.Sleep(lev.frameTimeout)
		screenService.ClearScreen()
	}
}

func processStaticObjects(
	staticObjects []services.ScreenObject,
	screenSvc *services.ScreenService,
) {
	for _, object := range staticObjects {
		screenSvc.Draw(object)
		object.Unblock()
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

func collectNewDeleteOld(
	objectChannel chan services.ScreenObject,
	screenService *services.ScreenService,
) ([][][]services.ScreenObject, []services.ScreenObject) {
Poll:
	for {
		select {
		case newObject := <-objectChannel:
			interactiveObjects[newObject] = true
		default:
			break Poll
		}
	}
	processedObjects := screenService.NewObjectList()
	invulnerableObjects := []services.ScreenObject{}

	for object := range interactiveObjects {
		if !object.IsActive() {
			delete(interactiveObjects, object)
			continue
		}
		if !object.IsVulnerable() {
			invulnerableObjects = append(invulnerableObjects, object)
			continue
		}
		coordinates, _ := object.GetViewCoordinates()
		for _, coord_pair := range coordinates {
			x, y := coord_pair[0], coord_pair[1]
			if screenService.IsInsideScreen(float64(x), float64(y)) {
				processedObjects[y][x] = append(processedObjects[y][x], object)
			}
		}
	}
	return processedObjects, invulnerableObjects
}

func processInteractiveObjects2(
	ctx context.Context,
	objectChannel chan services.ScreenObject,
	screenService *services.ScreenService,
	spaceshipCollisions, destroyedMeteorites int,
) (int, int) {
	active, passive := collectNewDeleteOld(objectChannel, screenService)
	for y, row := range active {
		for x, objects := range row {
			if len(objects) == 0 {
				continue
			}
			if len(objects) == 1 && !objects[0].GetDrawStatus() {
				screenService.Draw(objects[0])
				objects[0].MarkDrawn()
				active[y][x] = []services.ScreenObject{}
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
	for _, object := range passive {
		screenService.Draw(object)
	}

	for object := range interactiveObjects {
		if object.IsActive() {
			object.Unblock()
		}
	}
	return spaceshipCollisions, destroyedMeteorites
}
