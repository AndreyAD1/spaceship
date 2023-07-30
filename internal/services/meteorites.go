package services

import (
	"math/rand"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
)

const MeteoriteView = `  ___
 /   \
/     /
\____/`

var MeteoriteRuneView []rune = []rune{
	' ', ' ', '_', '_', '_', '\n',
	' ', '/', 0x85, 0x85, 0x85, '\\', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, 0x85, '/', '\n',
	'\\', '_', '_', '_', '_', '/',
}
var maxMeteoriteWidth = 7
var mutx sync.Mutex
var destroyedMeteorites = 0

func GenerateMeteorites(
	events chan ScreenObject,
	explosions chan ScreenObject,
	screenSvc *ScreenService,
	logger *log.Logger,
) {
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	width, _ := screenSvc.GetScreenSize()
	meteoritesOnUpperEdge := make([]*Meteorite, width)
Outer:	for {
		time.Sleep(time.Millisecond * 200)
		if rand.Float32() < 0.9 {
			continue
		}
		column := rand.Intn(width - 2)
		for i := column; i < column + maxMeteoriteWidth && i < width; i++ {
			if meteorite := meteoritesOnUpperEdge[i]; meteorite != nil {
				if meteorite.Y <= 1 && meteorite.Active {
					continue Outer
				}
				meteoritesOnUpperEdge[i] = nil
			}
		}
		baseObject := BaseObject{
			false,
			true,
			float64(column),
			-6,
			meteoriteStyle,
			0.02,
			string(MeteoriteRuneView),
			make(chan (struct{})),
			make(chan (struct{})),
		}
		meteorite := Meteorite{
			baseObject,
			events,
			screenSvc,
			explosions,
		}
		for i := column; i < column + maxMeteoriteWidth && i < width; i++ {
			meteoritesOnUpperEdge[i] = &meteorite
		}
		go meteorite.Move()
	}
}

type Meteorite struct {
	BaseObject
	Objects          chan<- ScreenObject
	ScreenSvc        ScreenSvc
	explosionChannel chan<- ScreenObject
}

func (meteorite *Meteorite) Move() {
	for {
		newY := meteorite.Y + meteorite.MaxSpeed
		_, height := meteorite.ScreenSvc.GetScreenSize()
		if newY > float64(height)+2 {
			meteorite.Active = false
			break
		}
		meteorite.Y = newY
		meteorite.Objects <- meteorite

		select {
		case <-meteorite.Cancel:
			meteorite.Active = false
			return
		case <-meteorite.UnblockCh:
		}
	}
}

func (meteorite *Meteorite) Collide(objects []ScreenObject) bool {
	allObjectsAreMeteorsOrSpaceship := true
	for _, obj := range objects {
		switch obj.(type) {
		case *Meteorite:
		case *Spaceship:
		default:
			allObjectsAreMeteorsOrSpaceship = false
		}
	}
	if !allObjectsAreMeteorsOrSpaceship {
		meteorite.Deactivate()
		go Explode(meteorite.explosionChannel, meteorite.X, meteorite.Y)
		mutx.Lock()
		defer mutx.Unlock()
		destroyedMeteorites++
		return true
	}
	return false
}
