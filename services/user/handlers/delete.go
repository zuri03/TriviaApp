package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type DeleteHandler struct {
	RedisHandler *rejson.Handler
	Signaler     chan os.Signal
}

func (d *DeleteHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		defer func() {
			d.Signaler <- os.Interrupt
		}()
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	_, err := d.RedisHandler.JSONDel(key, ".")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
		defer func() {
			d.Signaler <- os.Interrupt
		}()
		return
	}

	writer.Write([]byte("user deleted"))
}
