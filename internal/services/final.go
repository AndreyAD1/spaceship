package services

import "github.com/gdamore/tcell/v2"

const gameOverLabel = `
 _____                         ______               
/ ____|                       /  __  \                
| |  __  __ _ _ __ ___   ___  | |  | |_   _____ _ __ 
| | |_ |/ _| | '_ | _ \ / _ \ | |  | \ \ / / _ \ '__|
| |__| | (_| | | | | | |  __/ | |__| |\ V /  __/ |   
\ _____|\__,_|_| |_| |_|\___|  \____/  \_/ \___|_|   
`
const winLabel = `
__   _____  _   _  __        _____ _   _ _ 
\ \ / / _ \| | | | \ \      / /_ _| \ | | |
 \ V / | | | | | |  \ \ /\ / / | ||  \| | |
  | || |_| | |_| |   \ V  V /  | || |\  |_|
  |_| \___/ \___/     \_/\_/  |___|_| \_(_)
`
const nextLevel = `
 _  _         _     _                _       
| \| |_____ _| |_  | |   _____ _____| |      
| .' / -_) \ /  _| | |__/ -_) V / -_) |_ _ _ 
|_|\_\___/_\_\\__| |____\___|\_/\___|_(_|_|_)
`

func DrawGameOver(channel chan<- *BaseObject, screenSvc *ScreenService) {
	draw(channel, screenSvc, gameOverLabel)
}

func DrawWin(channel chan<- *BaseObject, screenSvc *ScreenService) {
	draw(channel, screenSvc, winLabel)
}

func DrawNextLevel(channel chan<- *BaseObject, screenSvc *ScreenService) {
	draw(channel, screenSvc, nextLevel)
}

func draw(channel chan<- *BaseObject, screenSvc *ScreenService, view string) {
	width, height := screenSvc.GetScreenSize()
	labelRow := width / 3
	labelColumn := height / 3
	gameover := BaseObject{
		false,
		true,
		float64(labelRow),
		float64(labelColumn),
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.01,
		view,
		make(chan struct{}),
		make(chan struct{}),
	}
	for {
		channel <- &gameover
	}
}
