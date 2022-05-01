package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type GetUserHandler struct {
	RedisHandler *rejson.Handler
}

func (c *CreateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("getting user request")

	if req.Header.Get("Content-Type") != "" {
		value := req.Header.Values("Content-Type")
		fmt.Printf("value => %s\n", value[0])
		if value[0] != "application/json" {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	fmt.Printf("getting user with key => %s\n", key)

	//For now just user username+password as a key
	userJson, err := c.RedisHandler.JSONGet(key, ".")
	if err != nil {
		fmt.Printf("error mess => %s\n", err.Error())
		if err.Error() != "" {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			return
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf("User %s does not exist => %s\n", userDetails.Username)))
			return
		}
	}

	fmt.Printf("found user of type => %T\n", userJson)
	writer.Write([]byte("found user"))
}
