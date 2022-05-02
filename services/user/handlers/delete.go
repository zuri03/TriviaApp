package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/models"
)

type DeleteHandler struct {
	RedisHandler *rejson.Handler
}

func (d *DeleteHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var userDetails models.UserDetails
	if err := json.NewDecoder(req.Body).Decode(&userDetails); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("error parsing json"))
		return
	}

	key := fmt.Sprintf("%s:%s", userDetails.Username, userDetails.Password)

	result, err := d.RedisHandler.JSONDel(key, ".")
	if err != nil {
		fmt.Printf("Error in delete => %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Error has occured => %s\n", err.Error())))
		return
	}

	fmt.Printf("type of deletion result %T\n", result)
	writer.Write([]byte("user deleted"))
}
