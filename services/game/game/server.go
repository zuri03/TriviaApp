package game

/*
TODO: Clean up code by seperating functions
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct{ j jservice }

type PlayerAnswer struct {
	Address string `json:"address"`
	Answer  string `json:"answer"`
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func (connection observer) updatePlayer(conn *websocket.Conn) {
	for {
		state := <-connection.messenger
		if err := conn.WriteMessage(websocket.TextMessage, state); err != nil {
			fmt.Printf("error in reader: %s\n", err)
			return
		}
	}
}

func initServer(g *game) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	//http.Handle("/Question", &Handler{j: jservice{Client: http.Client{}}})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, req *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true } //Remove this after development

		ws, err := upgrader.Upgrade(writer, req, nil)
		if err != nil {
			fmt.Printf("error with upgrader: %s\n", err)
		}
		o := observer{messenger: make(chan []byte)}
		g.RegisterObserver(o)
		o.updatePlayer(ws)
	})

	http.HandleFunc("/Submit", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		type response struct {
			Data  string `json:"data"`
			Error string `json:"error"`
		}
		if err := req.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			str, _ := json.Marshal(response{Data: "", Error: "Internal Server Error "})
			writer.Write(str)
			return
		}

		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			str, _ := json.Marshal(response{Data: "", Error: "Internal Server Error "})
			writer.Write(str)
			return
		}

		var playerAnswer PlayerAnswer
		if err := json.Unmarshal(body, &playerAnswer); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			str, _ := json.Marshal(response{Data: "", Error: "Internal Server Error "})
			writer.Write(str)
			return
		}

		if playerAnswer.Answer == g.currentState.Answer {
			result, err := g.bettingPotClient.declareWinner(playerAnswer.Address)
			if err != nil {
				fmt.Printf("betting service error => %s\n", err.Error())
			} else {
				fmt.Printf("results => %s\n", result)
			}
		}

		respObj := struct {
			Data string `json:"data"`
		}{
			Data: "Success",
		}
		jsonStr, _ := json.Marshal(respObj)
		writer.Write(jsonStr)
	})

	go func() {
		//err := http.ListenAndServe("localhost:8080", nil)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Printf("error occured: %s \n", err.Error())
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}
