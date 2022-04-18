package game

/**
GAME STATE IS CONTROLLED BY TWO FUNCTIONS CALLING EACH OTHER IN A LOOP
BOTH FUNCTIONS COMPLETE SOME ACTIONS THEN PAUSE FOR A CERTAIN AMOUNT OF TIME BEFORE CALLING THE
**/
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type phase int8

const (
	STOPPED       phase = 0
	BETTING             = 1
	PRE_ROUND           = 2
	ROUND_ONGOING       = 3
	POST_ROUND          = 4

	BETTTING_DURATION int = 3e10
	ROUND_DURATION    int = 4.5e10
)

//Take answer out of state
type state struct {
	Phase    phase  `json:"phase"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type game struct {
	observers        []observer
	jserviceClient   jservice
	bettingPotClient bettingPot
	currentState     state
}

type observer struct {
	messenger chan []byte
}

func (g *game) notifyObservers() {
	for _, o := range g.observers {
		jsonStr, _ := json.Marshal(g.currentState)
		o.messenger <- jsonStr
	}
}

func (g *game) RegisterObserver(o observer) {
	g.observers = append(g.observers, o)
}

/**
These functions only update the state and call their respective action function
*/
//Returns number of seconds the state lasts
func (g *game) startBetting() (phase, error) {
	g.currentState.Phase = BETTING
	g.notifyObservers()
	signaler := make(chan int8)
	defer close(signaler)
	go bettingActions(signaler)
	<-signaler
	return PRE_ROUND, nil //5 Seconds
}

//Pre round and post round will have to return channels to indicate when to change state
//So maybe wrap the ints
func (g *game) startPreRound() (phase, error) {
	g.currentState.Phase = PRE_ROUND
	g.notifyObservers()
	signaler := make(chan int8)
	defer close(signaler)
	go preGameActions(signaler, g)
	val := <-signaler
	if val == 0 {
		return STOPPED, errors.New("Error occured with jservice client")
	}
	return ROUND_ONGOING, nil
}

//Fires off new state event
//Returns next state and amount of time this state lasts for
//Update state variable
func (g *game) startRound() (phase, error) {
	g.currentState.Phase = ROUND_ONGOING
	g.notifyObservers()
	signaler := make(chan int8)
	defer close(signaler)
	go roundOngoingActions(signaler)
	<-signaler
	return POST_ROUND, nil
}

func (g *game) startPostRound() (phase, error) {
	g.currentState.Phase = POST_ROUND
	g.notifyObservers()
	signaler := make(chan int8)
	defer close(signaler)
	go postGameActions(signaler, g)
	<-signaler
	return BETTING, nil
}

func InitGame() {
	game := game{
		observers:      make([]observer, 0),
		jserviceClient: jservice{Client: http.Client{}},
		currentState: state{
			Phase:    BETTING,
			Question: "",
			Answer:   "",
		},
	}
	defer game.run() //may have to run this in a seperate go routine
	initServer(&game)
}

//Run function always assumes it is called when game has been stopped
//Run game logic and loop
func (g *game) run() {
	go func() {
		nextRound, err := g.startBetting()
		if err != nil {
			//handle error
		}
		fmt.Println("initial start betting returned")
		for {
			switch nextRound {
			case 1:
				nextRound, err = g.startBetting()
				break
			case 2:
				nextRound, err = g.startPreRound()
				break
			case 3:
				nextRound, err = g.startRound()
				break
			case 4:
				nextRound, err = g.startPostRound()
				break
			default:
				//Throw an error here
			}
			if err != nil {
				//Handle error
			}
			fmt.Println("FINISHED PHASE")
		}
	}()
}
