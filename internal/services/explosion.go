package services

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

const explosionFrameTimeout = time.Millisecond * 400

var explosionFrames []string = []string{
	`
     (_) 
 (  (   (  (
() (  (  )
  ( )  ()
    `,
	`
    (_) 
(  (   (   
  (  (  )
   )  (
    `,
	`
   (  
 (   (   
(     (
 )  (
    `,
	`
( 
  (
(  
    `,
}

func Explode(ch chan<- ScreenObject, XCentre, YCentre float64) {
	frameIndex := 0
	explosion := BaseObject{
		false,
		true,
		XCentre - 2,
		YCentre,
		tcell.StyleDefault.Background(tcell.ColorReset),
		maxSpeed,
		explosionFrames[frameIndex],
		make(chan struct{}),
		make(chan struct{}),
	}
	ticker := time.NewTicker(explosionFrameTimeout)
	defer ticker.Stop()
	for {
		ch <- &explosion
		select {
		case <-ticker.C:
			frameIndex++
			if frameIndex >= len(explosionFrames) {
				<-explosion.UnblockCh
				return
			}
			explosion.View = explosionFrames[frameIndex]
		default:
		}
		<-explosion.UnblockCh
	}
}