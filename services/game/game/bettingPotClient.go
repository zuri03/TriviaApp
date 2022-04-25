package game

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type bettingPot struct{ Client http.Client }

func (bettingPot *bettingPot) declareWinner(address string) (string, error) {

	data := url.Values{}
	data.Set("address", address)

	req, err := http.NewRequest("POST", "http://betting:8000/winner", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := bettingPot.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
