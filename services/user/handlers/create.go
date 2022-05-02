package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type CreateHandler struct {
	RedisHandler *rejson.Handler
}

func (c *CreateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("create user request recieved")

	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		//For now just return an internal server error
		writer.WriteHeader(http.StatusInternalServerError)
		//We will implement a standard json response for now just return plain text
		writer.Write([]byte("error parsing json"))
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)
	fmt.Printf("Key => %s\n", key)

	//For now just user username+password as a key
	_, err := c.RedisHandler.JSONGet(key, ".")
	if err != nil {
		fmt.Printf("error mess => %s\n", err.Error())

		if err.Error() != "redis: nil" {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			return
		}

		newUser := &models.User{
			Id:        uuid.New().String(),
			Username:  userDetails.Username,
			Password:  userDetails.Password,
			CreatedAt: time.Now().Format("02 Jan 2006 15:04:05"),
			Role:      "user",
			Wins:      0,
		}

		c.RedisHandler.JSONSet(key, ".", newUser)

		fmt.Println("Stored user json in redis")

		jsonBytes, _ := json.Marshal(newUser)

		writer.WriteHeader(http.StatusAccepted)
		writer.Write([]byte(jsonBytes))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte(fmt.Sprintf("user %s already exists", userDetails.Username)))
}
