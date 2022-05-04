package main

import (
	"fmt"

	"github.com/zuri03/game/game"
)

func main() {
	game := game.InitGame()
	<-game.Termination
	fmt.Println("SHUT DOWN SIGNAL RECEIVED")
}
