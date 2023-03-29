package services

import "github.com/gdamore/tcell/v2"

const label = `
_____                         ____                 
/ ____|                       / __ \                
| |  __  __ _ _ __ ___   ___  | |  | |_   _____ _ __ 
| | |_ |/ _| | '_ | _ \ / _ \ | |  | \ \ / / _ \ '__|
| |__| | (_| | | | | | |  __/ | |__| |\ V /  __/ |   
\_____|\__,_|_| |_| |_|\___|  \____/  \_/ \___|_|   
`

func DrawGameOver(channel chan<- *GameOver, screenSvc *ScreenService) {
	width, height := screenSvc.screen.Size()
	labelRow := width / 4
	labelColumn := height / 4
	base := BaseObject {
		false,
		true,
		float64(labelRow),
		float64(labelColumn),
		tcell.StyleDefault.Background(tcell.ColorReset),
		0.01,
	}
	gameover := GameOver{base, label}
	for {
		channel<- &gameover
	}
}

type GameOver struct {
	BaseObject
	view string
}