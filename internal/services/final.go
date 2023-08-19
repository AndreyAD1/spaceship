package services

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

const gameOverView string = `
 _____                         ______               
/ ____|                       /  __  \                
| |  __  __ _ _ __ ___   ___  | |  | |_   _____ _ __ 
| | |_ |/ _| | '_ | _ \ / _ \ | |  | \ \ / / _ \ '__|
| |__| | (_| | | | | | |  __/ | |__| |\ V /  __/ |   
\ _____|\__,_|_| |_| |_|\___|  \____/  \_/ \___|_|   
`
const winView string = `
__   _____  _   _  __        _____ _   _ _ 
\ \ / / _ \| | | | \ \      / /_ _| \ | | |
 \ V / | | | | | |  \ \ /\ / / | ||  \| | |
  | || |_| | |_| |   \ V  V /  | || |\  |_|
  |_| \___/ \___/     \_/\_/  |___|_| \_(_)
`
const nextLevelView string = `
 _  _         _     _                _       
| \| |_____ _| |_  | |   _____ _____| |      
| .' / -_) \ /  _| | |__/ -_) V / -_) |_ _ _ 
|_|\_\___/_\_\\__| |____\___|\_/\___|_(_|_|_)
`

type FinalLabel struct {
	view   string
	width  int
	height int
}

var GameOver = FinalLabel{gameOverView, 54, 8}
var Win = FinalLabel{winView, 44, 6}
var NextLevel = FinalLabel{nextLevelView, 46, 5}

func DrawLabel(
	ctx context.Context,
	channel chan<- *BaseObject,
	screenSvc *ScreenService,
	label FinalLabel,
) {
	width, height := screenSvc.GetScreenSize()
	labelRow := width/2 - label.width/2
	labelColumn := height/2 - label.height/2
	gameover := BaseObject{
		false,
		true,
		float64(labelRow),
		float64(labelColumn),
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.01,
		label.view,
		make(chan struct{}),
		make(chan struct{}),
	}
	for {
		select {
		case channel <- &gameover:
		case <-ctx.Done():
			return
		}
	}
}
