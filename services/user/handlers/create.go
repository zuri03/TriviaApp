package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type CreateHandler struct {
	RedisHandler *rejson.Handler
	Signaler     chan os.Signal
}

func (c *CreateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("create user request recieved")

	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		//For now just return an internal server error
		writer.WriteHeader(http.StatusInternalServerError)
		//We will implement a standard json response for now just return plain text
		writer.Write([]byte("error parsing json"))
		defer func() {
			c.Signaler <- os.Interrupt
		}()
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	//For now just user username+password as a key
	_, err := c.RedisHandler.JSONGet(key, ".")
	if err != nil {

		if err.Error() != "redis: nil" {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			defer func() {
				c.Signaler <- os.Interrupt
			}()
			return
		}

		newUser := &models.User{
			Id:        uuid.New().String(),
			Username:  userDetails.Username,
			Password:  userDetails.Password,
			CreatedAt: time.Now().Format("02 Jan 2006 15:04:05"),
			Role:      "user",
		}

		c.RedisHandler.JSONSet(key, ".", newUser)

		jsonBytes, _ := json.Marshal(newUser)

		fmt.Println("saved user")
		writer.WriteHeader(http.StatusAccepted)
		writer.Write([]byte(jsonBytes))
		return
	}

	fmt.Println("user already exist")
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte(fmt.Sprintf("user %s already exists", userDetails.Username)))
}
