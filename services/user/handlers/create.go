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

	//Check request header content type
	//Decode json in body
	//Ensure username and password does not already exist
	//create new user object with new id
	//store new user in redis

	fmt.Println("create user request recieved")

	if req.Header.Get("Content-Type") != "" {
		value := req.Header.Values("Content-Type")
		fmt.Printf("value => %s\n", value[0])
		if value[0] != "application/json" {
			writer.WriteHeader(http.StatusInternalServerError)
			//str, _ := json.Marshal(response{Data: "", Error: "Internal Server Error "})
			//writer.Write(str)
			return
		}
	}

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
		fmt.Printf("error == nil ? %b\n", err.Error() == "nil")
		fmt.Printf("error == empty string ? %b\n", err.Error() == "")
		/*
			if err.Error() != "" {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
				return
			}
		*/

		newUser := &models.User{
			Id:        uuid.New().String(),
			Username:  userDetails.Username,
			Password:  userDetails.Password,
			CreatedAt: time.Now().Format("02 Jan 2006 15:04:05"),
			Role:      "user",
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
