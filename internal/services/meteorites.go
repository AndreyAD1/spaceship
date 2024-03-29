package services

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const MeteoriteView1 = `  ____
 /    \
/     /
\____/`
const MeteoriteView2 = `  ___
 _/  \
/  __/
\_/`
const MeteoriteView3 = `  _______
/        \
\         \
 \_       /
   \_____/`

var MeteoriteRuneView1 []rune = []rune{
	' ', ' ', '_', '_', '_', '_', '\n',
	' ', '/', 0x85, 0x85, 0x85, 0x85, '\\', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, 0x85, '/', '\n',
	'\\', '_', '_', '_', '_', '/',
}
var MeteoriteRuneView2 []rune = []rune{
	' ', ' ', ' ', '_', '_', '\n',
	' ', '_', '/', 0x85, 0x85, '\\', '\n',
	'/', 0x85, 0x85, '_', '_', '/', '\n',
	'\\', '_', '/', '\n',
}
var MeteoriteRuneView3 []rune = []rune{
	' ', '_', '_', '_', '_', '_', '_', '_', '_', '\n',
	'/', 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, '\\', '\n',
	'\\', 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, '\\', '\n',
	' ', '\\', '_', 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, 0x85, '/', '\n',
	' ', ' ', ' ', '\\', '_', '_', '_', '_', '_', '/', '\n',
}

type meteoriteProps struct {
	view     []rune
	maxWidth int
}

var meteoriteMutx sync.RWMutex
var destroyedMeteorites = 0
var meteorites = []meteoriteProps{
	{MeteoriteRuneView1, 7},
	{MeteoriteRuneView2, 6},
	{MeteoriteRuneView3, 10},
}

func GenerateMeteorites(
	ctx context.Context,
	events chan ScreenObject,
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
		meteoriteType := meteorites[rand.Intn(len(meteorites))]
		for i := column; i < column+meteoriteType.maxWidth && i < width; i++ {
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
			string(meteoriteType.view),
			make(chan (struct{})),
			make(chan (struct{})),
			true,
		}
		meteorite := Meteorite{
			baseObject,
			events,
			screenSvc,
		}
		for i := column; i < column+meteoriteType.maxWidth && i < width; i++ {
			meteoritesOnUpperEdge[i] = &meteorite
		}
		go meteorite.Move(ctx)
	}
}

type Meteorite struct {
	BaseObject
	Objects   chan<- ScreenObject
	ScreenSvc ScreenSvc
}

func (meteorite *Meteorite) Move(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case meteorite.Objects <- meteorite:
	}
	for {
		newY := meteorite.Y + meteorite.MaxSpeed
		_, height := meteorite.ScreenSvc.GetScreenSize()
		if newY > float64(height)+2 {
			meteorite.Active = false
			break
		}
		meteorite.Y = newY

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
		go Explode(ctx, meteorite.Objects, meteorite.X, meteorite.Y)
		meteoriteMutx.Lock()
		defer meteoriteMutx.Unlock()
		destroyedMeteorites++
		return true
	}
	return false
}
