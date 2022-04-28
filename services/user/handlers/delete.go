package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Username struct {
	Username string `json:"username"`
}

type DeleteHandler struct{}

func (c *DeleteHandler) ServeHttp(writer http.ResponseWriter, req *http.Request) {
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
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		str, _ := json.Marshal(response{Data: "", Error: "Internal Server Error "})
		writer.Write(str)
		return
	}
}
