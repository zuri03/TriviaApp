package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type GetHandler struct {
	RedisHandler *rejson.Handler
	Signaler     chan os.Signal
}

func (g *GetHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	fmt.Println("getting user request")

	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		defer func() {
			g.Signaler <- os.Interrupt
		}()
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	//For now just user username+password as a key
	userJson, err := g.RedisHandler.JSONGet(key, ".")
	if err != nil {
		if err.Error() == "redis: nil" {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf("User %s does not exist\n", userDetails.Username)))
			return
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
			defer func() {
				g.Signaler <- os.Interrupt
			}()
			return
		}
	}

	var user models.User
	if err := json.Unmarshal(userJson.([]byte), &user); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Error converting json to obj => %s\n", err.Error())))
		defer func() {
			g.Signaler <- os.Interrupt
		}()
		return
	}

	writer.Write(userJson.([]byte))
}
