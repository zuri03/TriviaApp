package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/handlers"
)

func InitRouter(redisJson *rejson.Handler, signaler chan os.Signal) *http.ServeMux {

	router := http.NewServeMux()

	fmt.Printf("signaler == nil? %t\n", signaler == nil)
	createHandler := handlers.CreateHandler{RedisHandler: redisJson, Signaler: signaler}
	getHandler := handlers.GetHandler{RedisHandler: redisJson, Signaler: signaler}
	deleteHandler := handlers.DeleteHandler{RedisHandler: redisJson, Signaler: signaler}
	updateHandler := handlers.UpdateHandler{RedisHandler: redisJson, Signaler: signaler}

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "" {
			value := req.Header.Values("Content-Type")
			if value[0] != "application/json" {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Error: expected json"))
				return
			}
		}

		fmt.Println("inside of user router")
		switch req.Method {
		case http.MethodGet:
			getHandler.ServeHTTP(writer, req)
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
			deleteHandler.ServeHTTP(writer, req)
		case http.MethodPut:
			updateHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	return router
}
