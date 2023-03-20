package main

import (
	"fmt"
	"time"

	"github.com/AndreyAD1/spaceship/frame"
)

func main() {
	fmt.Println("a spaceship game will be here")
	newFrame := frame.DrawFrame()
	defer newFrame.Clear()
	time.Sleep(10 * time.Second)
}