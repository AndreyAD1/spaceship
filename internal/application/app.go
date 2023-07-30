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

func (app Application) Run(is_last_level bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = log.WithContext(ctx, app.Logger)
	screenService, err := services.NewScreenService()
	if err != nil {
		return err
	}
	defer screenService.Finish()

	levelConfigs := []levelConfig {
		{5, 3, false},
		{7, 2, true},
	}
	
	for _, levelConfig := range levelConfigs {
		level := NewLevel(levelConfig, app.FrameTimeout)
		err = level.Run(ctx, screenService)
	}
	return err
}
