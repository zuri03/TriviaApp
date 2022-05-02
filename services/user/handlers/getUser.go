package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type GetHandler struct {
	RedisHandler *rejson.Handler
}

func (c *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("getting user request")

	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	//For now just user username+password as a key
	userJson, err := c.RedisHandler.JSONGet(key, ".")
	if err != nil {
		fmt.Printf("error mess => %s\n", err.Error())
		fmt.Printf("error redsi:nil => %t\n", err.Error() == "redis: nil")
		if err.Error() == "redis: nil" {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf("User %s does not exist\n", userDetails.Username)))
			return
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			return
		}
	}

	fmt.Printf("found user of type => %T\n", userJson)

	var user models.User
	if err := json.Unmarshal(userJson.([]byte), &user); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Error converting json to obj => %s\n", err.Error())))
		return
	}
	fmt.Printf("got user => \n %+v\n", user)
	writer.Write(userJson.([]byte))
}
