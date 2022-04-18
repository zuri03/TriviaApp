package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strings"
)

type jserviceResponse struct {
	Random []randomResponse `json:"random"`
}
type randomResponse struct {
	Id           int      `json:"id"`
	Answer       string   `json:"answer"`
	Question     string   `json:"question"`
	Value        int      `json:"value"`
	Airdate      string   `json:"airdate"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	CategoryId   int      `json:"category_id"`
	GameId       int      `json:"game_id"`
	InvalidCount int      `json:"invalid_count"`
	Category     category `json:"category"`
}

type category struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	CluesCount int    `json:"clues_count"`
}

type jservice struct{ Client http.Client }

func (jservice *jservice) getQuestionAndAnswer() (jserviceResponse, error) {
	var response jserviceResponse
	req, err := http.NewRequest("GET", "https://jservice.io/api/random", nil)
	if err != nil {
		return response, err
	}

	resp, err := jservice.Client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsonStr := fmt.Sprintf("{\"random\":%s}", string(body))
	json.Unmarshal([]byte(jsonStr), &response)

	return response, nil
}
