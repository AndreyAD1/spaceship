package services

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const MeteoriteView1 = `  ___
 /   \
/     /
\____/`
const MeteoriteView2 = ` ____
/    \
\____/`
 const MeteoriteView3 = `  __
 /  \____
/        \
\         )
 \_______/`

var MeteoriteRuneView1 []rune = []rune{
	' ', ' ', '_', '_', '_', '\n',
	' ', '/', 0x85, 0x85, 0x85, '\\', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, 0x85, '/', '\n',
	'\\', '_', '_', '_', '_', '/',
}
var MeteoriteRuneView2 []rune = []rune{
	' ', '_', '_', '_', '_', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, '\\', '\n',
	'\\', '_', '_', '_', '_', '/', '\n',
}
var MeteoriteRuneView3 []rune = []rune{
	' ', ' ', '_', '_', '\n',
	' ', '/', 0x85, 0x85, '\\', '_', '_', '_', '_', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, '\\', '\n',
	'\\', 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, ')', '\n',
	' ', '\\', '_', '_', '_', '_', '_', '_', '_', '/', '\n',
}
var maxMeteoriteWidth = 7
var mutx sync.Mutex
var destroyedMeteorites = 0
var views = [][]rune{MeteoriteRuneView1, MeteoriteRuneView2, MeteoriteRuneView3}

func GenerateMeteorites(
	ctx context.Context,
	events chan ScreenObject,
	explosions chan ScreenObject,
	screenSvc *ScreenService,
) {
	destroyedMeteorites = 0
	meteoriteStyle := tcell.StyleDefault.Background(tcell.ColorReset)
	width, _ := screenSvc.GetScreenSize()
	meteoritesOnUpperEdge := make([]*Meteorite, width)
Outer:
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		time.Sleep(time.Millisecond * 200)
		if rand.Float32() < 0.9 {
			continue
		}
		column := rand.Intn(width - 2)
		for i := column; i < column+maxMeteoriteWidth && i < width; i++ {
			if meteorite := meteoritesOnUpperEdge[i]; meteorite != nil {
				if meteorite.Y <= 1 && meteorite.Active {
					continue Outer
				}
				meteoritesOnUpperEdge[i] = nil
			}
		}
		meteoriteView := views[rand.Intn(len(views))]
		baseObject := BaseObject{
			false,
			true,
			float64(column),
			-6,
			meteoriteStyle,
			0.02,
			string(meteoriteView),
			make(chan (struct{})),
			make(chan (struct{})),
		}
		meteorite := Meteorite{
			baseObject,
			events,
			screenSvc,
			explosions,
		}
		for i := column; i < column+maxMeteoriteWidth && i < width; i++ {
			meteoritesOnUpperEdge[i] = &meteorite
		}
		go meteorite.Move(ctx)
	}
}

type Meteorite struct {
	BaseObject
	Objects          chan<- ScreenObject
	ScreenSvc        ScreenSvc
	explosionChannel chan<- ScreenObject
}

func (meteorite *Meteorite) Move(ctx context.Context) {
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
		case <-ctx.Done():
			return
		case <-meteorite.Cancel:
			meteorite.Active = false
			return
		case <-meteorite.UnblockCh:
		}
	}
}

func (meteorite *Meteorite) Collide(ctx context.Context, objects []ScreenObject) bool {
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
		go Explode(ctx, meteorite.explosionChannel, meteorite.X, meteorite.Y)
		mutx.Lock()
		defer mutx.Unlock()
		destroyedMeteorites++
		return true
	}
	return false
}
