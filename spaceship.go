package main

import (
	"fmt"

	"github.com/AndreyAD1/spaceship/frame"
)

func main() {
	fmt.Println("a spaceship game will be here")
	newFrame := frame.DrawFrame()
	fmt.Println(newFrame)
}