package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type UpdateHandler struct {
	RedisHandler *rejson.Handler
	Signaler     chan os.Signal
}

func (u *UpdateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("create user request recieved")

	var updatedUser models.User
	if err := json.NewDecoder(req.Body).Decode(&updatedUser); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		signal.Notify(u.Signaler, os.Interrupt)
		return
	}

	if updatedUser.Role == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Missing new role"))
		return
	}

	key := fmt.Sprintf("%s:%s", updatedUser.Username, updatedUser.Password)
	fmt.Printf("Key => %s\n", key)

	//For now just user username+password as a key
	oldUserJson, err := u.RedisHandler.JSONGet(key, ".")
	if err != nil {
		fmt.Printf("error mess => %s\n", err.Error())
		fmt.Printf("error redsi:nil => %t\n", err.Error() == "redis: nil")
		if err.Error() == "redis: nil" {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf("User %s does not exist\n", updatedUser.Username)))
			return
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			defer func() {
				u.Signaler <- os.Interrupt
			}()
			return
		}
	}

	var oldUser models.User
	_ = json.Unmarshal(oldUserJson.([]byte), &oldUser)
	temp := updatedUser.Role
	updatedUser = oldUser
	updatedUser.Role = temp

	u.RedisHandler.JSONSet(key, ".", updatedUser)
	writer.Write([]byte("user updated"))
}
