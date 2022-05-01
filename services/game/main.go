package main

import (
	"fmt"

	"github.com/zuri03/game/game"
)

func main() {
	game := game.InitGame()
	val := <-game.Termination
	if val == 0 {
		fmt.Printf("Graceful Shutdown")
	} else {
		fmt.Printf("Error has occured shutting down")
	}
}
