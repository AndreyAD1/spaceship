package application

import (
	"context"
	"runtime"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = log.WithContext(ctx, app.Logger)
	screenService, err := services.NewScreenService()
	if err != nil {
		return err
	}
	defer screenService.Finish()

	levelConfigs := []levelConfig{
		{"Level 1", 5, 3, false},
		{"Level 2", 7, 2, true},
	}

	for i, levelConfig := range levelConfigs {
		level := NewLevel(levelConfig, app.FrameTimeout)
		err = level.Run(ctx, screenService)
		app.Logger.Debugf(
			"goroutines after a level %v: %v",
			i,
			runtime.NumGoroutine(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
