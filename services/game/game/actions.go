package game

import (
	"fmt"
	"time"
)

/*
This package implements all of the actions the game must take while it is in a specific phase
*/

//All of the post game logic goes here
//
func postGameActions(signaler chan int8, g *game) {
	g.currentState.Answer = ""
	g.currentState.Question = ""
	time.Sleep(time.Duration(5000000000))
	signaler <- 1
	return
}

func preGameActions(signaler chan int8, g *game) {
	res, err := g.jserviceClient.getQuestionAndAnswer()
	if err != nil {
		signaler <- 0
	}
	g.currentState.Question = res.Random[0].Question
	g.currentState.Answer = res.Random[0].Answer
	signaler <- 1
	return
}

func bettingActions(signaler chan int8) {
	fmt.Printf("starting betting \n")
	time.Sleep(time.Duration(BETTTING_DURATION)) //Replace after development
	signaler <- 1
	return
}

func roundOngoingActions(signaler chan int8) {
	fmt.Printf("starting round \n")
	time.Sleep(time.Duration(ROUND_DURATION))
	signaler <- 1
	return
}
